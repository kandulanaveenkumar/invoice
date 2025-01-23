package workflow

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IFlowInstanceParams interface {
	Upsert(ctx *context.Context, m *models.FlowInstanceParams) error
	Get(ctx *context.Context, id string) (*models.FlowInstanceParams, error)
	GetAll(ctx *context.Context) ([]*models.FlowInstanceParams, error)
	Delete(ctx *context.Context, id string) error
	GetForFlowInstance(ctx *context.Context, flowInstanceID string) ([]*models.FlowInstanceParams, error)
}

type FlowInstanceParams struct {
}

func NewFlowInstanceParams() IFlowInstanceParams {
	return &FlowInstanceParams{}
}

func (t *FlowInstanceParams) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".flow_instance_params"
}

func (t *FlowInstanceParams) Upsert(ctx *context.Context, m *models.FlowInstanceParams) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *FlowInstanceParams) Get(ctx *context.Context, id string) (*models.FlowInstanceParams, error) {
	var result models.FlowInstanceParams
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *FlowInstanceParams) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.FlowInstanceParams{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *FlowInstanceParams) GetForFlowInstance(ctx *context.Context, flowInstanceID string) ([]*models.FlowInstanceParams, error) {
	var result []*models.FlowInstanceParams
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("flow_instance_id = ?", flowInstanceID).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *FlowInstanceParams) GetAll(ctx *context.Context) ([]*models.FlowInstanceParams, error) {
	var result []*models.FlowInstanceParams
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}
