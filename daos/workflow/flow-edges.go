package workflow

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IFlowEdges interface {
	Upsert(ctx *context.Context, m *models.FlowEdges) error
	Get(ctx *context.Context, id string) (*models.FlowEdges, error)
	GetAll(ctx *context.Context) ([]*models.FlowEdges, error)
	Delete(ctx *context.Context, id string) error
	GetForWorkflow(ctx *context.Context, workflowID string) ([]*models.FlowEdges, error)
}

type FlowEdges struct {
}

func NewFlowEdges() IFlowEdges {
	return &FlowEdges{}
}

func (t *FlowEdges) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".flow_edges"
}

func (t *FlowEdges) Upsert(ctx *context.Context, m *models.FlowEdges) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *FlowEdges) Get(ctx *context.Context, id string) (*models.FlowEdges, error) {
	var result models.FlowEdges
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *FlowEdges) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.FlowEdges{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *FlowEdges) GetForWorkflow(ctx *context.Context, workflowID string) ([]*models.FlowEdges, error) {
	var result []*models.FlowEdges
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("workflow_id = ?", workflowID).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *FlowEdges) GetAll(ctx *context.Context) ([]*models.FlowEdges, error) {
	var result []*models.FlowEdges
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}
