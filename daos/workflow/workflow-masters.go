package workflow

import (
	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IWorkflowMasters interface {
	Upsert(ctx *context.Context, m *models.WorkflowMasters) error
	Get(ctx *context.Context, id string) (*models.WorkflowMasters, error)
	GetAll(ctx *context.Context) ([]*models.WorkflowMasters, error)
	Delete(ctx *context.Context, id string) error
	GetForFilterCount(ctx *context.Context, req *dtos.GetFilter) (int, error)
	GetForFilter(ctx *context.Context, req *dtos.GetFilter) ([]*models.WorkflowMasters, error)
	GetForFilterWithoutPagination(ctx *context.Context, req *dtos.GetFilter) ([]*models.WorkflowMasters, error)
}

type WorkflowMasters struct {
}

func NewWorkflowMasters() IWorkflowMasters {
	return &WorkflowMasters{}
}

func (t *WorkflowMasters) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".workflow_masters"
}

func (t *WorkflowMasters) Upsert(ctx *context.Context, m *models.WorkflowMasters) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *WorkflowMasters) Get(ctx *context.Context, id string) (*models.WorkflowMasters, error) {
	var result models.WorkflowMasters
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *WorkflowMasters) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.WorkflowMasters{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *WorkflowMasters) GetAll(ctx *context.Context) ([]*models.WorkflowMasters, error) {
	var result []*models.WorkflowMasters
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *WorkflowMasters) GetForFilterCount(ctx *context.Context, req *dtos.GetFilter) (int, error) {
	var cnt int
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("count(*)")
	if req.Q != "" {
		tx = tx.Where("name ILIKE ?", req.Q+"%")
	}

	if req.ID != "" {
		tx = tx.Where("id = ?", req.ID)
	}

	if req.Type != "" {
		tx = tx.Where("type = ?", req.Type)
	}

	if req.RegionID != "" {
		tx = tx.Where("region_id = ?", req.RegionID)
	}

	err := tx.Scan(&cnt).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return 0, err
	}

	return cnt, err
}

func (t *WorkflowMasters) GetForFilter(ctx *context.Context, req *dtos.GetFilter) ([]*models.WorkflowMasters, error) {
	var result []*models.WorkflowMasters
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if req.Q != "" {
		tx = tx.Where("name ILIKE ?", req.Q+"%")
	}

	if req.ID != "" {
		tx = tx.Where("id = ?", req.ID)
	}

	if req.Type != "" {
		tx = tx.Where("type = ?", req.Type)
	}

	if req.RegionID != "" {
		tx = tx.Where("region_id = ?", req.RegionID)
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

func (t *WorkflowMasters) GetForFilterWithoutPagination(ctx *context.Context, req *dtos.GetFilter) ([]*models.WorkflowMasters, error) {
	var result []*models.WorkflowMasters
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if req.Q != "" {
		tx = tx.Where("name ILIKE ?", req.Q+"%")
	}

	if req.ID != "" {
		tx = tx.Where("id = ?", req.ID)
	}

	if req.Type != "" {
		tx = tx.Where("type = ?", req.Type)
	}

	if req.RegionID != "" {
		tx = tx.Where("region_id = ?", req.RegionID)
	}

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return result, err
}
