package masterconsignee

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IMasterConsignee interface {
	Upsert(ctx *context.Context, m ...*models.MasterConsignee) error
	Get(ctx *context.Context, id string) (*models.MasterConsignee, error)
	GetAll(ctx *context.Context, ids []string, shipmentId string) ([]*models.MasterConsignee, error)
	Delete(ctx *context.Context, id string) error
}

type MasterConsignee struct {
}

func NewMasterConsignee() IMasterConsignee {
	return &MasterConsignee{}
}

func (t *MasterConsignee) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "master_consignees"
}

func (t *MasterConsignee) Upsert(ctx *context.Context, m ...*models.MasterConsignee) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *MasterConsignee) Get(ctx *context.Context, id string) (*models.MasterConsignee, error) {
	var result models.MasterConsignee
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get masterconsignee.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *MasterConsignee) Delete(ctx *context.Context, id string) error {
	var result models.MasterConsignee
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete masterconsignee.", zap.Error(err))
		return err
	}

	return err
}

func (t *MasterConsignee) GetAll(ctx *context.Context, ids []string, shipmentId string) ([]*models.MasterConsignee, error) {
	var result []*models.MasterConsignee

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if len(ids) > 0 {
		tx.Where("id IN ?", ids)
	}

	if shipmentId != "" {
		tx.Where("shipment_id = ?", shipmentId)
	}

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get masterconsignees.", zap.Error(err))
		return nil, err
	}

	return result, err
}
