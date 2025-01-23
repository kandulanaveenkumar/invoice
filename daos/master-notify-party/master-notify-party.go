package masternotifyparty

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IMasterNotifyParty interface {
	Upsert(ctx *context.Context, m ...*models.MasterNotifyParty) error
	Get(ctx *context.Context, id string) (*models.MasterNotifyParty, error)
	GetAll(ctx *context.Context, ids []string, shipmentId string) ([]*models.MasterNotifyParty, error)
	Delete(ctx *context.Context, id string) error
}

type MasterNotifyParty struct {
}

func NewMasterNotifyParty() IMasterNotifyParty {
	return &MasterNotifyParty{}
}

func (t *MasterNotifyParty) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "master_notify_parties"
}

func (t *MasterNotifyParty) Upsert(ctx *context.Context, m ...*models.MasterNotifyParty) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *MasterNotifyParty) Get(ctx *context.Context, id string) (*models.MasterNotifyParty, error) {
	var result models.MasterNotifyParty
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get masternotifyparty.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *MasterNotifyParty) Delete(ctx *context.Context, id string) error {
	var result models.MasterNotifyParty
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete masternotifyparty.", zap.Error(err))
		return err
	}

	return err
}

func (t *MasterNotifyParty) GetAll(ctx *context.Context, ids []string, shipmentId string) ([]*models.MasterNotifyParty, error) {
	var result []*models.MasterNotifyParty

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if len(ids) > 0 {
		tx.Where("id IN ?", ids)
	}

	if shipmentId != "" {
		tx.Where("shipment_id = ?", shipmentId)
	}

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get masternotifypartys.", zap.Error(err))
		return nil, err
	}

	return result, err
}
