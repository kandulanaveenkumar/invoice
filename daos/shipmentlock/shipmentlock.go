package shipmentlock

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IShipmentLock interface {
	Upsert(ctx *context.Context, m ...*models.ShipmentLock) error
	Get(ctx *context.Context, id string) (*models.ShipmentLock, error)
	GetAll(ctx *context.Context) ([]*models.ShipmentLock, error)

	GetByShipment(ctx *context.Context, shipmentId uuid.UUID) ([]*models.ShipmentLock, error)
	UpdateShipmentLock(ctx *context.Context, IsBookingLocked bool, Id uuid.UUID) error
}

type ShipmentLock struct {
}

func NewShipmentLock() IShipmentLock {
	return &ShipmentLock{}
}

func (t *ShipmentLock) getTable(ctx *context.Context) string {
	if ctx.TenantID == "" {
		ctx.Log.Error("TenantID is empty")
		return "public.shipment_lock" 
	}
	return ctx.TenantID + "." + "shipment_lock"
}

func (t *ShipmentLock) Upsert(ctx *context.Context, m ...*models.ShipmentLock) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *ShipmentLock) Get(ctx *context.Context, id string) (*models.ShipmentLock, error) {
	var result models.ShipmentLock
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get ShipmentLock.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *ShipmentLock) GetAll(ctx *context.Context) ([]*models.ShipmentLock, error) {
	var result []*models.ShipmentLock

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Where("is_locked = ?", false).
		Find(&result).Error

	if err != nil {
		ctx.Log.Error("Unable to get shipment locks", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (t *ShipmentLock) GetByShipment(ctx *context.Context, shipmentId uuid.UUID) ([]*models.ShipmentLock, error) {
	var result []*models.ShipmentLock
	err := ctx.DB.Table(t.getTable(ctx)).Where("shipment_id = ?", shipmentId).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get ShipmentLocks.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *ShipmentLock) UpdateShipmentLock(ctx *context.Context, IsBookingLocked bool, Id uuid.UUID) error {
	err := ctx.DB.Debug().
		Table(t.getTable(ctx)).
		Where("id = ?", Id).
		Update("is_locked", IsBookingLocked).Error

	if err != nil {
		return err
	}

	return nil
}
