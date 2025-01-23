package mastermultiplehbl

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IMasterMultipleHbl interface {
	Upsert(ctx *context.Context, m ...*models.MasterMultipleHbl) error
	Get(ctx *context.Context, shipmentId string) (*models.MasterMultipleHbl, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.MasterMultipleHbl, error)
	Delete(ctx *context.Context, id string) error
}

type MasterMultipleHbl struct {
}

func NewMasterMultipleHbl() IMasterMultipleHbl {
	return &MasterMultipleHbl{}
}

func (t *MasterMultipleHbl) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "master_multiple_hbls"
}

func (t *MasterMultipleHbl) Upsert(ctx *context.Context, m ...*models.MasterMultipleHbl) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *MasterMultipleHbl) Get(ctx *context.Context, shipmentId string) (*models.MasterMultipleHbl, error) {
	var result models.MasterMultipleHbl
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "shipment_id = ?", shipmentId).Error
	if err != nil {
		ctx.Log.Error("Unable to get mastermultiplehbl.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *MasterMultipleHbl) Delete(ctx *context.Context, id string) error {
	var result models.MasterMultipleHbl
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete mastermultiplehbl.", zap.Error(err))
		return err
	}

	return err
}

func (t *MasterMultipleHbl) GetAll(ctx *context.Context, ids []string) ([]*models.MasterMultipleHbl, error) {
	var result []*models.MasterMultipleHbl
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get mastermultiplehbls.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get mastermultiplehbls.", zap.Error(err))
		return nil, err
	}

	return result, err
}
