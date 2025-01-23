package shipmentregionrestriction

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IShipmentRegionRestriction interface {
	Upsert(ctx *context.Context, m ...*models.ShipmentRegionRestriction) error
	Get(ctx *context.Context, id string) (*models.ShipmentRegionRestriction, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.ShipmentRegionRestriction, error)
	Delete(ctx *context.Context, id string) error
}

type ShipmentRegionRestriction struct {
}

func NewShipmentRegionRestriction() IShipmentRegionRestriction {
	return &ShipmentRegionRestriction{}
}

func (t *ShipmentRegionRestriction) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "shipment_region_restrictions"
}

func (t *ShipmentRegionRestriction) Upsert(ctx *context.Context, m ...*models.ShipmentRegionRestriction) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *ShipmentRegionRestriction) Get(ctx *context.Context, id string) (*models.ShipmentRegionRestriction, error) {
	var result models.ShipmentRegionRestriction
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentregionrestriction.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *ShipmentRegionRestriction) Delete(ctx *context.Context, id string) error {
	var result models.ShipmentRegionRestriction
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete shipmentregionrestriction.", zap.Error(err))
		return err
	}

	return err
}

func (t *ShipmentRegionRestriction) GetAll(ctx *context.Context, ids []string) ([]*models.ShipmentRegionRestriction, error) {
	var result []*models.ShipmentRegionRestriction
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get shipmentregionrestrictions.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentregionrestrictions.", zap.Error(err))
		return nil, err
	}

	return result, err
}
