package workflow

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IFlowInstanceEdges interface {
	Upsert(ctx *context.Context, m *models.FlowInstanceEdges) error
	Get(ctx *context.Context, id string) (*models.FlowInstanceEdges, error)
	GetAll(ctx *context.Context) ([]*models.FlowInstanceEdges, error)
	Delete(ctx *context.Context, id string) error
	GetUsingFromFlowID(ctx *context.Context, flowID string) ([]*models.FlowInstanceEdges, error)
	GetForInstance(ctx *context.Context, instanceID string) ([]*models.FlowInstanceEdges, error)
}

type FlowInstanceEdges struct {
}

func NewFlowInstanceEdges() IFlowInstanceEdges {
	return &FlowInstanceEdges{}
}

func (t *FlowInstanceEdges) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".flow_instance_edges"
}

func (t *FlowInstanceEdges) Upsert(ctx *context.Context, m *models.FlowInstanceEdges) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *FlowInstanceEdges) Get(ctx *context.Context, id string) (*models.FlowInstanceEdges, error) {
	var result models.FlowInstanceEdges
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *FlowInstanceEdges) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.FlowInstanceEdges{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *FlowInstanceEdges) GetUsingFromFlowID(ctx *context.Context, flowID string) ([]*models.FlowInstanceEdges, error) {
	var result []*models.FlowInstanceEdges
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("from_flow_id = ?", flowID).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *FlowInstanceEdges) GetForInstance(ctx *context.Context, instanceID string) ([]*models.FlowInstanceEdges, error) {
	var result []*models.FlowInstanceEdges
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("instance_id = ?", instanceID)

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *FlowInstanceEdges) GetAll(ctx *context.Context) ([]*models.FlowInstanceEdges, error) {
	var result []*models.FlowInstanceEdges
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}
