package workflow

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IFlows interface {
	Upsert(ctx *context.Context, m *models.Flows) error
	Get(ctx *context.Context, id string) (*models.Flows, error)
	GetAll(ctx *context.Context) ([]*models.Flows, error)
	Delete(ctx *context.Context, id string) error
	GetForWorkflow(ctx *context.Context, workflowID string) ([]*models.Flows, error)
}

type Flows struct {
}

func NewFlows() IFlows {
	return &Flows{}
}

func (t *Flows) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".flows"
}

func (t *Flows) Upsert(ctx *context.Context, m *models.Flows) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *Flows) Get(ctx *context.Context, id string) (*models.Flows, error) {
	var result models.Flows
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *Flows) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.Flows{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *Flows) GetForWorkflow(ctx *context.Context, workflowID string) ([]*models.Flows, error) {
	var result []*models.Flows
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("workflow_id = ?", workflowID).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Flows) GetAll(ctx *context.Context) ([]*models.Flows, error) {
	var result []*models.Flows
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}
