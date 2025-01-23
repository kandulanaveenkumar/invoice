package workflow

import (
	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IServices interface {
	Upsert(ctx *context.Context, m *models.Services) error
	Get(ctx *context.Context, id string) (*models.Services, error)
	GetAll(ctx *context.Context) ([]*models.Services, error)
	Delete(ctx *context.Context, id string) error
	GetForFilterCount(ctx *context.Context, req *dtos.GetFilter) (int, error)
	GetForFilter(ctx *context.Context, req *dtos.GetFilter) ([]*models.Services, error)
	GetForFilterWithoutPagination(ctx *context.Context, req *dtos.GetFilter) ([]*models.Services, error)
}

type Services struct {
}

func NewServices() IServices {
	return &Services{}
}

func (t *Services) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".services"
}

func (t *Services) Upsert(ctx *context.Context, m *models.Services) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *Services) Get(ctx *context.Context, id string) (*models.Services, error) {
	var result models.Services
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *Services) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.Services{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *Services) GetAll(ctx *context.Context) ([]*models.Services, error) {
	var result []*models.Services
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Services) GetForFilterCount(ctx *context.Context, req *dtos.GetFilter) (int, error) {
	var cnt int
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("count(*)")
	if req.Q != "" {
		tx = tx.Where("name ILIKE ?", req.Q+"%")
	}

	if req.ID != "" {
		tx = tx.Where("id = ?", req.ID)
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

func (t *Services) GetForFilter(ctx *context.Context, req *dtos.GetFilter) ([]*models.Services, error) {
	var result []*models.Services
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if req.Q != "" {
		tx = tx.Where("name ILIKE ?", req.Q+"%")
	}

	if req.ID != "" {
		tx = tx.Where("id = ?", req.ID)
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

func (t *Services) GetForFilterWithoutPagination(ctx *context.Context, req *dtos.GetFilter) ([]*models.Services, error) {
	var result []*models.Services
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if req.Q != "" {
		tx = tx.Where("name ILIKE ?", req.Q+"%")
	}

	if req.ID != "" {
		tx = tx.Where("id = ?", req.ID)
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
