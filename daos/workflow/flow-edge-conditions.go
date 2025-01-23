package workflow

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IFlowEdgeConditions interface {
	Upsert(ctx *context.Context, m *models.FlowEdgeConditions) error
	Get(ctx *context.Context, id string) (*models.FlowEdgeConditions, error)
	GetAll(ctx *context.Context) ([]*models.FlowEdgeConditions, error)
	Delete(ctx *context.Context, id string) error
	GetForFlowEdge(ctx *context.Context, edgeID string) ([]*models.FlowEdgeConditions, error)
}

type FlowEdgeConditions struct {
}

func NewFlowEdgeConditions() IFlowEdgeConditions {
	return &FlowEdgeConditions{}
}

func (t *FlowEdgeConditions) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".flow_edge_conditions"
}

func (t *FlowEdgeConditions) Upsert(ctx *context.Context, m *models.FlowEdgeConditions) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *FlowEdgeConditions) Get(ctx *context.Context, id string) (*models.FlowEdgeConditions, error) {
	var result models.FlowEdgeConditions
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *FlowEdgeConditions) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.FlowEdgeConditions{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *FlowEdgeConditions) GetForFlowEdge(ctx *context.Context, edgeID string) ([]*models.FlowEdgeConditions, error) {
	var result []*models.FlowEdgeConditions
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("flow_edge_id = ?", edgeID).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *FlowEdgeConditions) GetAll(ctx *context.Context) ([]*models.FlowEdgeConditions, error) {
	var result []*models.FlowEdgeConditions
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}
