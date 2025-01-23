package shipmentlocation

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IShipmentLocation interface {
	Upsert(ctx *context.Context, m ...*models.ShipmentLocation) error
	UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.ShipmentLocation) error

	Get(ctx *context.Context, shipmentId string, locationType string) (*models.ShipmentLocation, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.ShipmentLocation, error)
	Delete(ctx *context.Context, id string) error
	DeleteByShipmentId(ctx *context.Context, shipmentId string) error
}

type ShipmentLocation struct {
}

func NewShipmentLocation() IShipmentLocation {
	return &ShipmentLocation{}
}

func (t *ShipmentLocation) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "shipment_locations"
}

func (t *ShipmentLocation) Upsert(ctx *context.Context, m ...*models.ShipmentLocation) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *ShipmentLocation) UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.ShipmentLocation) error {
	return tx.Table(t.getTable(ctx)).Save(m).Error
}

func (t *ShipmentLocation) Get(ctx *context.Context, shipmentId string, locationType string) (*models.ShipmentLocation, error) {
	var result models.ShipmentLocation
	err := ctx.DB.WithContext(ctx.Request.Context()).
		Table(t.getTable(ctx)).
		Where("shipment_id = ? AND type = ?", shipmentId, locationType).
		First(&result).
		Error
	return &result, err
}

func (t *ShipmentLocation) Delete(ctx *context.Context, id string) error {
	var result models.ShipmentLocation
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete shipmentlocation.", zap.Error(err))
		return err
	}

	return err
}

func (t *ShipmentLocation) DeleteByShipmentId(ctx *context.Context, shipmentId string) error {
	var result models.ShipmentLocation
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "shipment_id = ?", shipmentId).Error
	if err != nil {
		ctx.Log.Error("Unable to delete shipmentlocations.", zap.Error(err))
		return err
	}

	return err
}

func (t *ShipmentLocation) GetAll(ctx *context.Context, ids []string) ([]*models.ShipmentLocation, error) {
	var result []*models.ShipmentLocation
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get shipmentlocations.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("shipment_id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentlocations.", zap.Error(err))
		return nil, err
	}

	return result, err
}
