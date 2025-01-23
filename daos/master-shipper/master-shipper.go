package mastershipper

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IMasterShipper interface {
	Upsert(ctx *context.Context, m ...*models.MasterShipper) error
	Get(ctx *context.Context, id string) (*models.MasterShipper, error)
	GetAll(ctx *context.Context, ids []string, shipmentId string) ([]*models.MasterShipper, error)
	Delete(ctx *context.Context, id string) error
}

type MasterShipper struct {
}

func NewMasterShipper() IMasterShipper {
	return &MasterShipper{}
}

func (t *MasterShipper) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "master_shippers"
}

func (t *MasterShipper) Upsert(ctx *context.Context, m ...*models.MasterShipper) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *MasterShipper) Get(ctx *context.Context, id string) (*models.MasterShipper, error) {
	var result models.MasterShipper
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get mastershipper.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *MasterShipper) Delete(ctx *context.Context, id string) error {
	var result models.MasterShipper
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete mastershipper.", zap.Error(err))
		return err
	}

	return err
}

func (t *MasterShipper) GetAll(ctx *context.Context, ids []string, shipmentId string) ([]*models.MasterShipper, error) {
	var result []*models.MasterShipper

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if len(ids) > 0 {
		tx.Where("id IN ?", ids)
	}

	if shipmentId != "" {
		tx.Where("shipment_id = ?", shipmentId)
	}

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get mastershippers.", zap.Error(err))
		return nil, err
	}

	return result, err
}
