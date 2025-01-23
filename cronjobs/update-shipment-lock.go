package cronjobs

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/apis/id"
	"bitbucket.org/radarventures/forwarder-adapters/apis/misc"
	"bitbucket.org/radarventures/forwarder-adapters/apis/workflow"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	cards "bitbucket.org/radarventures/forwarder-shipments/daos/card"
	"bitbucket.org/radarventures/forwarder-shipments/daos/lineitem"
	"bitbucket.org/radarventures/forwarder-shipments/daos/rfq"
	"bitbucket.org/radarventures/forwarder-shipments/daos/shipment"
	"bitbucket.org/radarventures/forwarder-shipments/daos/shipmentlock"
	cardsAssignment "bitbucket.org/radarventures/forwarder-shipments/services/card-assignment"
	shipments "bitbucket.org/radarventures/forwarder-shipments/services/shipment"
	"github.com/google/uuid"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"bitbucket.org/radarventures/forwarder-shipments/services/card"
	rfqSer "bitbucket.org/radarventures/forwarder-shipments/services/rfq"

	"go.uber.org/zap"
)

type UpdateShipmentLock struct {
	id         id.ID
	cardDb     cards.ICard
	shipmentDb shipment.IShipment
	rfqDb      rfq.IRfq
	misc       misc.Misc
	rfqSer     rfqSer.IRfqService
	lineitem   lineitem.ILineItem
	workflow   *workflow.Workflow
	card       card.ICardService
	cards      cardsAssignment.ICardAssignmentService

	shipment     shipments.IShipmentService
	shipmentlock shipmentlock.IShipmentLock
}

func NewUpdateShipmentLock() IShipmentLock {
	return &UpdateShipmentLock{
		id:           *id.New(config.Get().IdURL),
		cardDb:       cards.NewCard(),
		shipmentDb:   shipment.NewShipment(),
		rfqDb:        rfq.NewRfq(),
		misc:         *misc.New(config.Get().MiscURL),
		rfqSer:       rfqSer.NewRfqService(),
		workflow:     workflow.New(config.Get().WorkflowURL),
		lineitem:     lineitem.NewLineItem(),
		card:         card.NewCardService(),
		shipment:     shipments.NewShipmentService(),
		shipmentlock: shipmentlock.NewShipmentLock(),
		cards:        cardsAssignment.NewCardAssignmentService(),
	}
}

type IShipmentLock interface {
	UpdateShipmentlockStatus(ctx *context.Context)
	UpdateShipmentlockStatusV2(ctx *context.Context)
}

func (c *UpdateShipmentLock) UpdateShipmentlockStatus(ctx *context.Context) {

	ctx.Log.Info("inside UpdateShipmentLock")

	shipments, err := c.shipmentDb.GetAllUnlockedShipments(ctx)
	if err != nil {
		ctx.Log.Error("Unable to get the shipments", zap.Error(err))
		return
	}

	for _, shipment := range shipments {
		if shipment.Type == globals.BookingTypeCONSOL {
			continue
		}

		ctx.Log.Info("shipment is", zap.Any("shipment id", shipment.Id))
		milestones, err := c.workflow.GetMilestones(ctx, shipment.Id.String())
		if err != nil {
			ctx.Log.Error("unable to get the milestones", zap.Error(err))
			ctx.Log.Info("EF-1354 milestoneDate", zap.Any("shipment", shipment.Id))
			continue
		}

		lineItems, err := c.lineitem.GetLineItemsWithFilter(ctx, &models.LiFields{
			QuoteId: shipment.QuoteId,
		})
		if err != nil {
			ctx.Log.Error("error getting line items", zap.Error(err))
			continue
		}

		milestoneCheckForLockShipment := false
		milestoneDate := ""
		foundBookingConfirmed := false

		if milestones != nil && milestones.Milestones != nil {
			for _, milestone := range milestones.Milestones {
				milestoneName := strings.ToLower(milestone.Name)
				normalizedMilestoneName := strings.ToLower(strings.TrimSpace(milestoneName))

				if strings.Contains(normalizedMilestoneName, "booking confirm") {
					if milestone.Status == globals.StatusCompleted && milestone.CompletedAt != 0 {
						foundBookingConfirmed = true
						milestoneDate = time.Unix(milestone.CompletedAt, 0).String()
					}
				}
			}

			if !foundBookingConfirmed {
				for _, milestone := range milestones.Milestones {
					if milestone.Name == constants.MilestoneBookingCreated && milestone.Status == globals.StatusCompleted && milestone.CompletedAt != 0 {
						// Handle timestamp precision (milliseconds or seconds)
						if milestone.CompletedAt > 9999999999 {
							milestoneDate = time.Unix(milestone.CompletedAt/1000, 0).String() // Assuming milliseconds
						} else {
							milestoneDate = time.Unix(milestone.CompletedAt, 0).String() // Assuming seconds
						}
						break
					}
				}
			}
			ctx.Log.Info("milestone date", zap.Any("milestonedate", milestoneDate), zap.Any("shipment_id", shipment.Id))

			if milestoneDate != "" {
				parsedTime, err := parseMilestoneDate(milestoneDate)
				if err != nil {
					ctx.Log.Error("Error parsing milestoneDate: ", zap.Error(err))
					continue
				}

				if time.Since(parsedTime).Hours()/24 > 60 && (shipment.Type == globals.BookingTypeFCL || shipment.Type == globals.BookingTypeLCL) {
					milestoneCheckForLockShipment = true
				} else if time.Since(parsedTime).Hours()/24 > 30 && (shipment.Type == globals.BookingTypeAIR || shipment.Type == constants.ShipmentTypeMisc) {
					milestoneCheckForLockShipment = true
				}
				ctx.Log.Info("parsed time", zap.Any("parsed time", time.Since(parsedTime).Hours()), zap.Any("shipment_id", shipment.Id))
			}

		}
		lineItemscheckforLockShipment := true
		for i := range lineItems {
			if (lineItems[i].IsSellInvoiceGenerated == nil || !*lineItems[i].IsSellInvoiceGenerated) ||
				(lineItems[i].IsBuyApproved == nil || !*lineItems[i].IsBuyApproved) {
				lineItemscheckforLockShipment = false
				break
			}
		}
		ctx.Log.Info("entered outside", zap.Any("milestoneCheckForLockShipment", milestoneCheckForLockShipment), zap.Any("shipment id", shipment.Id))
		if milestoneCheckForLockShipment || lineItemscheckforLockShipment {
			ctx.Log.Info("entered milestoneCheckForLockShipment", zap.Any("milestoneCheckForLockShipment", milestoneCheckForLockShipment))
			shipment.IsShipmentLocked = true
			err = c.shipmentDb.UpdateShipmentLock(ctx, shipment.Id.String(), shipment.IsShipmentLocked)
			if err != nil {
				ctx.Log.Error("error while updating shipment", zap.Any("shipment_id", shipment.Id), zap.Error(err))
				continue
			}

			c.cards.DeleteAllcards(ctx, &dtos.CardRequest{
				InstanceId:  shipment.Id.String(),
				ExecutiveId: config.Get().WizBotID,
			},
			)
			shipmentLock := &models.ShipmentLock{
				Id:         uuid.New(),
				ShipmentId: shipment.Id,
				UpdatedBy:  config.Get().WizBotID,
				IsLocked:   true,
			}

			err = c.shipment.UpdateShipmentLockStatusTimeline(ctx, shipmentLock, true, config.Get().WizBotID, "", false, shipment.Type, shipment.RegionId.String())
			if err != nil {
				ctx.Log.Error("error while persisting audit", zap.Error(err))
				continue
			}
			ctx.Log.Info("completed lock shipment", zap.Any("shipment is", shipment.Id), zap.Any("shipmentLock", shipmentLock))
		}
	}
	ctx.Log.Info("migration finished for updating shipment lock status")
}

func (c *UpdateShipmentLock) UpdateShipmentlockStatusV2(ctx *context.Context) {

	allShipmentLock, err := c.shipmentlock.GetAll(ctx)
	if err != nil {
		ctx.Log.Error("error while fetching unlocked shipments", zap.Error(err))
		return
	}

	for _, shipmentlock := range allShipmentLock {
		ctx.Log.Info("shipmentlockV2", zap.Any("shipmentlock.ShipmentId", shipmentlock.ShipmentId.String()))
		shipment, err := c.shipmentDb.Get(ctx, shipmentlock.ShipmentId.String())
		if err != nil {
			ctx.Log.Error("Unable to get shipment", zap.Error(err))
			return
		}

		if shipment.IsShipmentLocked {
			continue
		}

		if time.Since(shipment.UpdatedAt) > 90*time.Minute {
			shipment.IsShipmentLocked = true

			err = c.shipmentDb.UpdateShipmentLock(ctx, shipment.Id.String(), shipment.IsShipmentLocked)
			if err != nil {
				ctx.Log.Error("error while updating shipment", zap.Any("shipment_id", shipment.Id), zap.Error(err))
				continue
			}

			shipmentLock := &models.ShipmentLock{
				Id:         uuid.New(),
				ShipmentId: shipment.Id,
				UpdatedBy:  config.Get().WizBotID,
				IsLocked:   true,
			}
			ctx.Log.Info("EF-1354 inside cron2", zap.Any("regionId", shipment.Id))

			err = c.shipment.UpdateShipmentLockStatusTimeline(ctx, shipmentLock, true, config.Get().WizBotID, shipmentlock.Id.String(), true, shipment.Type, shipment.RegionId.String())
			if err != nil {
				ctx.Log.Error("error while persisting audit", zap.Error(err))
				continue
			}
		}

	}

	ctx.Log.Info("migration finished for updating shipment status V2")
}

func parseMilestoneDate(milestoneDate string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",                    // Date format
		"2006-01-02 15:04:05 -0700 MST", // Full timestamp with timezone
	}

	var parsedTime time.Time
	var err error
	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, milestoneDate)
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", milestoneDate)
}
