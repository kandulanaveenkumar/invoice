package workflow

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IFlowInstanceEdgeConditions interface {
	Upsert(ctx *context.Context, m *models.FlowInstanceEdgeConditions) error
	Get(ctx *context.Context, id string) (*models.FlowInstanceEdgeConditions, error)
	GetAll(ctx *context.Context) ([]*models.FlowInstanceEdgeConditions, error)
	Delete(ctx *context.Context, id string) error
	GetForEdge(ctx *context.Context, edgeID string) ([]*models.FlowInstanceEdgeConditions, error)
}

type FlowInstanceEdgeConditions struct {
}

func NewFlowInstanceEdgeConditions() IFlowInstanceEdgeConditions {
	return &FlowInstanceEdgeConditions{}
}

func (t *FlowInstanceEdgeConditions) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".flow_instance_edge_conditions"
}

func (t *FlowInstanceEdgeConditions) Upsert(ctx *context.Context, m *models.FlowInstanceEdgeConditions) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *FlowInstanceEdgeConditions) Get(ctx *context.Context, id string) (*models.FlowInstanceEdgeConditions, error) {
	var result models.FlowInstanceEdgeConditions
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *FlowInstanceEdgeConditions) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.FlowInstanceEdgeConditions{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *FlowInstanceEdgeConditions) GetForEdge(ctx *context.Context, edgeID string) ([]*models.FlowInstanceEdgeConditions, error) {
	var result []*models.FlowInstanceEdgeConditions
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("flow_edge_id = ?", edgeID).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *FlowInstanceEdgeConditions) GetAll(ctx *context.Context) ([]*models.FlowInstanceEdgeConditions, error) {
	var result []*models.FlowInstanceEdgeConditions
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}
