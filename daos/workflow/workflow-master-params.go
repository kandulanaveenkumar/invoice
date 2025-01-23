package workflow

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IWorkflowMasterParams interface {
	Upsert(ctx *context.Context, m *models.WorkflowMasterParams) error
	Get(ctx *context.Context, id string) (*models.WorkflowMasterParams, error)
	GetAll(ctx *context.Context) ([]*models.WorkflowMasterParams, error)
	Delete(ctx *context.Context, id string) error
	GetForWorkflowMaster(ctx *context.Context, masterID string) ([]*models.WorkflowMasterParams, error)
}

type WorkflowMasterParams struct {
}

func NewWorkflowMasterParams() IWorkflowMasterParams {
	return &WorkflowMasterParams{}
}

func (t *WorkflowMasterParams) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".workflow_master_params"
}

func (t *WorkflowMasterParams) Upsert(ctx *context.Context, m *models.WorkflowMasterParams) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *WorkflowMasterParams) Get(ctx *context.Context, id string) (*models.WorkflowMasterParams, error) {
	var result models.WorkflowMasterParams
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *WorkflowMasterParams) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.WorkflowMasterParams{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *WorkflowMasterParams) GetForWorkflowMaster(ctx *context.Context, masterID string) ([]*models.WorkflowMasterParams, error) {
	var result []*models.WorkflowMasterParams
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("workflow_master_id = ?", masterID).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *WorkflowMasterParams) GetAll(ctx *context.Context) ([]*models.WorkflowMasterParams, error) {
	var result []*models.WorkflowMasterParams
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}
