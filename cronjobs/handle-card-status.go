package cronjobs

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/apis/id"
	"bitbucket.org/radarventures/forwarder-adapters/apis/misc"
	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/misc"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/daos/card"
	cardaudits "bitbucket.org/radarventures/forwarder-shipments/daos/card-audits"
	"bitbucket.org/radarventures/forwarder-shipments/daos/rfq"
	"bitbucket.org/radarventures/forwarder-shipments/daos/shipment"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"bitbucket.org/radarventures/forwarder-shipments/services/websocket"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"go.uber.org/zap"
)

type HandleCardStatus struct {
	id           id.ID
	cardDb       card.ICard
	cardAuditsDb cardaudits.ICardAudits
	shipmentDb   shipment.IShipment
	rfqDb        rfq.IRfq
	misc         misc.Misc
	ws           websocket.WebSocket
}

func NewHandleCardStatus() IHandleCardStatus {
	return &HandleCardStatus{
		id:           *id.New(config.Get().IdURL),
		cardDb:       card.NewCard(),
		cardAuditsDb: cardaudits.NewCardAudits(),
		shipmentDb:   shipment.NewShipment(),
		rfqDb:        rfq.NewRfq(),
		misc:         *misc.New(config.Get().MiscURL),
		ws:           *websocket.NewWebSocket(),
	}
}

type IHandleCardStatus interface {
	HandleStatus(ctx *context.Context) error
}

// HandleStatus updates the status of pending cards based on their estimate time.
// If the current time is after the estimate time, the card status is set to "Breached",
// and a goroutine is started to update in the booking service.
// If the current time is within 60 minutes before the estimate time and the card is not already
// in a warning or breached status, the card status is set to "Warning".
// If the current time is before the estimate time and the card is already in a breached status,
// the card status is reset to "Created".

func (j *HandleCardStatus) HandleStatus(ctx *context.Context) error {

	ctx.Log.Info("HandleStatus Job Started")

	cards, err := j.cardDb.GetAllPendingCards(ctx)
	if err != nil {
		ctx.Log.Error("error while getting pending cards", zap.Error(err))
		return err
	}

	now := time.Now().UTC()

	cardStatus := ""
	var assignedToIds []string

	for _, card := range cards {

		cardStatus = card.Status

		// If the current time is after the estimate and the status is not breached
		if now.After(card.Estimate) && card.Status != constants.CardStatusBreached {
			card.Status = constants.CardStatusBreached

			if len(card.EscalatedById) > 0 {
				escID := card.EscalatedById[len(card.EscalatedById)-1]
				assignedToIds = append(assignedToIds, escID)

				if len(card.EscalatedById) > 1 {
					managerID := card.EscalatedById[len(card.EscalatedById)-2]
					assignedToIds = append(assignedToIds, managerID)
				}
			}

			assignedToIds = utils.AppendWithoutDuplicates(assignedToIds, card.AssignedTo)

		} else if now.After(card.Estimate.Add(-60*time.Minute)) && card.Status != constants.CardStatusWarning && card.Status != constants.CardStatusBreached {

			// If the current time is after the estimate minus minutes,
			// and the status is neither breached nor warning

			card.Status = constants.CardStatusWarning

			j.chatGenerationForWarning(ctx, &card)

			if len(card.EscalatedById) > 0 {
				escID := card.EscalatedById[len(card.EscalatedById)-1]
				assignedToIds = append(assignedToIds, escID)

				if len(card.EscalatedById) > 1 {
					managerID := card.EscalatedById[len(card.EscalatedById)-2]
					assignedToIds = append(assignedToIds, managerID)
				}
			}
			assignedToIds = utils.AppendWithoutDuplicates(assignedToIds, card.AssignedTo)

		} else if (card.Status == constants.CardStatusBreached && now.Before(card.Estimate)) || (card.Status == constants.CardStatusWarning && now.Before(card.Estimate.Add(-60*time.Minute))) {

			// Otherwise, if the current time is before the estimate and the status is breached,
			// reset the status to "Created"

			card.Status = constants.CardStatusCreated

			if len(card.EscalatedById) > 0 {
				escID := card.EscalatedById[len(card.EscalatedById)-1]
				assignedToIds = append(assignedToIds, escID)

				if len(card.EscalatedById) > 1 {
					managerID := card.EscalatedById[len(card.EscalatedById)-2]
					assignedToIds = append(assignedToIds, managerID)
				}
			}

			assignedToIds = utils.AppendWithoutDuplicates(assignedToIds, card.AssignedTo)
		}

		// Update the card status in the database
		if card.Status != cardStatus {

			reasonmap := make(map[string]interface{})
			reasonmap["card name"] = card.Name
			reasonmap["old card status"] = cardStatus
			reasonmap["new card status"] = card.Status

			cardAudit := &models.CardAudits{
				CardId:         card.Id,
				Name:           card.Name,
				InstanceId:     card.InstanceId,
				InstanceType:   card.InstanceType,
				Department:     card.Department,
				Status:         card.Status,
				FlowInstanceId: card.FlowInstanceId.String(),
				Reason:         reasonmap,
			}

			j.cardAuditsDb.Upsert(ctx, cardAudit)

			err = j.cardDb.UpdateStatus(ctx, card.Id.String(), card.Status)
			if err != nil {
				ctx.Log.Error("error updating card", zap.Error(err))
				return err
			}
		}

		// //Websocket message
		j.ws.SendCardsDataMiddleware(ctx, &card, map[string]interface{}{
			"card_id":     card.Id,
			"event":       constants.CardActionUpdate,
			"assigned_to": card.AssignedTo,
		}, assignedToIds)
	}

	ctx.Log.Info("HandleStatus Job Ended")

	return nil
}

// chatGenerationForWarning generates a chat message to warn about task delay for a given card
func (j *HandleCardStatus) chatGenerationForWarning(ctx *context.Context, card *models.Card) {

	executive, err := j.id.GetAccountInternal(ctx, card.AssignedTo)
	if err != nil {
		ctx.Log.Error("failed to get account details", zap.Error(err), zap.Any("aid", card.AssignedTo))
		return
	}

	var refId string
	var refType string
	var refCode string
	isShipment := false
	isRfq := false

	if card.InstanceId != "" {

		refId = card.InstanceId

		if card.InstanceType == constants.WorkflowTypeShipment {

			refType = constants.WorkflowTypeShipment
			isShipment = true

			shipment, err := j.shipmentDb.Get(ctx, refId)
			if err != nil {
				ctx.Log.Error("failed to get shipment details", zap.Error(err), zap.Any("sid", card.InstanceId))
				return
			}

			refCode = shipment.Code

		} else if card.InstanceType == constants.WorkflowTypeRFQ {

			refType = constants.WorkflowTypeRFQ
			isRfq = true

			rfq, err := j.rfqDb.Get(ctx, refId)
			if err != nil {
				ctx.Log.Error("failed to get rfq details", zap.Error(err), zap.Any("rfqid", card.InstanceId))
				return
			}

			refCode = rfq.Code
		}

	}

	// If the card is neither related to a rfq nor a shipment, no need to generate a chat message
	if !isRfq && isShipment {
		return
	}

	// Create a map to store tagged members for the chat message.
	taggedMembers := make(map[string]interface{})
	taggedMembers[executive.ReportingManager.String()] = executive.ReportingManager.String()
	taggedMembers[card.AssignedTo] = executive.Name

	reportingManager, err := j.id.GetAccountInternal(ctx, executive.ReportingManager.String())
	if err != nil {
		ctx.Log.Error("failed to get reporting manager details", zap.Error(err), zap.Any("rmid", executive.ReportingManager))
		return
	}

	var reportingManagerName string

	if reportingManager != nil {
		reportingManagerName = reportingManager.Name
	}

	// Generate the TaskLink with appropriate placeholders
	link := fmt.Sprintf("%s/dashboard/executive/%s?showback=true&&name=%s&&back=gv&open_task_id=%s&booking_code=%s", config.Get().BaseURL, card.AssignedTo, strings.Replace(executive.Name, " ", "%20", -1), card.Id, refCode)

	_, err = j.misc.SendCollab(ctx, &dtos.CollabMsg{
		RefID:         refId,
		RefType:       refType,
		Msg:           fmt.Sprintf("@%s The %s task is currently delayed and nearing its expiration deadline. ###Click*here~%s~### to take immediate action.\n \n cc: @%s", executive.Name, card.Name, link, reportingManagerName),
		TaggedMembers: taggedMembers,
		TaskRegionID:  card.RegionId,
		ChatType:      "internal_chat",
		AccountId:     config.Get().WizBotID,
	})
	if err != nil {
		ctx.Log.Error("failed to send collab message", zap.Error(err))
		return
	}

}
