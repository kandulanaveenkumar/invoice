package shipmentparty

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type IShipmentParty interface {
	Upsert(ctx *context.Context, m ...*models.ShipmentParty) error
	GetAll(ctx *context.Context, shipmentId []string, partyTypes []string) ([]*models.ShipmentParty, error)
	GetShipmentAddressDetailsForDSR(ctx *context.Context, shipmentIds []string) ([]models.ShipmentAddressDetail, error)
	GetShipmentListSearch(ctx *context.Context) ([]*models.ShipmentSearchFilter, error)
}

type ShipmentParty struct {
}

func NewShipmentParty() IShipmentParty {
	return &ShipmentParty{}
}

func (t *ShipmentParty) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "shipment_parties"
}

func (t *ShipmentParty) Upsert(ctx *context.Context, m ...*models.ShipmentParty) error {
	return ctx.DB.Debug().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "shipment_id"}, {Name: "address_id"}, {Name: "type"}},
		DoUpdates: clause.AssignmentColumns([]string{"region_id", "is_mpb", "created_at", "created_by", "updated_at", "updated_by"}),
	}).Table(t.getTable(ctx)).Save(m).Error
}

func (t *ShipmentParty) GetAll(ctx *context.Context, shipmentId []string, partyTypes []string) ([]*models.ShipmentParty, error) {
	var result []*models.ShipmentParty

	tx := ctx.DB.Debug().Debug().Table(t.getTable(ctx))

	if len(partyTypes) > 0 {
		tx.Where("type IN ?", partyTypes)
	}

	if len(shipmentId) > 0 {
		tx.Where("shipment_id IN ?", shipmentId)
	}

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentpartys.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *ShipmentParty) GetShipmentAddressDetailsForDSR(ctx *context.Context, shipmentIds []string) ([]models.ShipmentAddressDetail, error) {
	var results []models.ShipmentAddressDetail
	err := ctx.DB.Debug().Debug().Table(t.getTable(ctx)).
		Where("shipment_id IN (?) ", shipmentIds).
		Select("shipment_id, address_id").
		Scan(&results).Error
	if err != nil {
		ctx.Log.Error("error while getting shipment address details", zap.Error(err))
		return nil, err
	}

	return results, nil
}

// Will not work as name column is not present
func (t *ShipmentParty) GetShipmentListSearch(ctx *context.Context) ([]*models.ShipmentSearchFilter, error) {
	var shipmentSearch []*models.ShipmentSearchFilter
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("shipment_parties.id", "shipment_parties.name", "shipment_parties.type", "shipment_parties.created_at").
		Joins("JOIN "+ctx.TenantID+".shipments ON shipments.id = shipment_parties.shipment_id").
		Where("is_deleted = false AND ? = ANY(ARRAY[shipments.region_id, shipments.origin_region_id, shipments.dest_region_id])", ctx.Account.RegionID).
		Where("shipment_parties.type IN (?)", []string{constants.ShipmentContactShipper, constants.ShipmentContactConsignee}).
		Where("shipment_parties.name ilike ?", ctx.Query("q")+"%").
		Scan(&shipmentSearch).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentcontainer.", zap.Error(err))
		return nil, err
	}
	return shipmentSearch, nil
}
