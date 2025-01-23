package workflow

import (
	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IWorkflowFields interface {
	Upsert(ctx *context.Context, m *models.WorkflowFields) error
	Get(ctx *context.Context, id string) (*models.WorkflowFields, error)
	GetAll(ctx *context.Context) ([]*models.WorkflowFields, error)
	Delete(ctx *context.Context, id string) error
	GetForFilterCount(ctx *context.Context, req *dtos.GetFilter) (int, error)
	GetForFilter(ctx *context.Context, req *dtos.GetFilter) ([]*models.WorkflowFields, error)
}

type WorkflowFields struct {
}

func NewWorkflowFields() IWorkflowFields {
	return &WorkflowFields{}
}

func (t *WorkflowFields) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".workflow_fields"
}

func (t *WorkflowFields) Upsert(ctx *context.Context, m *models.WorkflowFields) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *WorkflowFields) Get(ctx *context.Context, id string) (*models.WorkflowFields, error) {
	var result models.WorkflowFields
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *WorkflowFields) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.WorkflowFields{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *WorkflowFields) GetAll(ctx *context.Context) ([]*models.WorkflowFields, error) {
	var result []*models.WorkflowFields
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *WorkflowFields) GetForFilterCount(ctx *context.Context, req *dtos.GetFilter) (int, error) {
	var cnt int
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("count(*)")
	if req.Q != "" {
		tx = tx.Where("field_name ILIKE ?", req.Q+"%")
	}

	if req.ID != "" {
		tx = tx.Where("id = ?", req.ID)
	}

	err := tx.Scan(&cnt).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return 0, err
	}

	return cnt, err
}

func (t *WorkflowFields) GetForFilter(ctx *context.Context, req *dtos.GetFilter) ([]*models.WorkflowFields, error) {
	var result []*models.WorkflowFields
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if req.Q != "" {
		tx = tx.Where("field_name ILIKE ?", req.Q+"%")
	}

	if req.ID != "" {
		tx = tx.Where("id = ?", req.ID)
	}

	tx = tx.Offset(req.Offset)
	tx = tx.Limit(req.Limit)

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return result, err
}
