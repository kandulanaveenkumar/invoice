package shipmentproduct

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IShipmentProduct interface {
	Upsert(ctx *context.Context, m ...*models.ShipmentProduct) error
	UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.ShipmentProduct) error

	Get(ctx *context.Context, id string) (*models.ShipmentProduct, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.ShipmentProduct, error)
	Delete(ctx *context.Context, id string) error
	GetByShipment(ctx *context.Context, shipmentIds []string) ([]*models.ShipmentProduct, error)
	DeleteByShipmentId(ctx *context.Context, shipmentId string) error
}

type ShipmentProduct struct {
}

func NewShipmentProduct() IShipmentProduct {
	return &ShipmentProduct{}
}

func (t *ShipmentProduct) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "shipment_products"
}

func (t *ShipmentProduct) Upsert(ctx *context.Context, m ...*models.ShipmentProduct) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *ShipmentProduct) UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.ShipmentProduct) error {
	return tx.Table(t.getTable(ctx)).Save(m).Error
}

func (t *ShipmentProduct) Get(ctx *context.Context, id string) (*models.ShipmentProduct, error) {
	var result models.ShipmentProduct
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentproduct.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *ShipmentProduct) Delete(ctx *context.Context, id string) error {
	var result models.ShipmentProduct
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete shipmentproduct.", zap.Error(err))
		return err
	}

	return err
}

func (t *ShipmentProduct) DeleteByShipmentId(ctx *context.Context, shipmentId string) error {
	var result models.ShipmentProduct
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "shipment_id = ?", shipmentId).Error
	if err != nil {
		ctx.Log.Error("Unable to delete shipmentproducts.", zap.Error(err))
		return err
	}

	return err
}

func (t *ShipmentProduct) GetAll(ctx *context.Context, ids []string) ([]*models.ShipmentProduct, error) {
	var result []*models.ShipmentProduct
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get shipmentproducts.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("shipment_id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentproducts.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *ShipmentProduct) GetByShipment(ctx *context.Context, shipmentIds []string) ([]*models.ShipmentProduct, error) {
	var result []*models.ShipmentProduct
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result, "shipment_id IN (?) AND is_deleted = false", shipmentIds).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipment products", zap.Error(err))
		return nil, err
	}

	return result, err
}
