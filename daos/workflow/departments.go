package workflow

import (
	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IDepartments interface {
	Upsert(ctx *context.Context, m *models.Departments) error
	Get(ctx *context.Context, id string) (*models.Departments, error)
	GetAll(ctx *context.Context) ([]*models.Departments, error)
	Delete(ctx *context.Context, id string) error
	GetForFilterCount(ctx *context.Context, req *dtos.GetFilter) (int, error)
	GetForFilter(ctx *context.Context, req *dtos.GetFilter) ([]*models.Departments, error)
	GetForRegionId(ctx *context.Context, regionId string) ([]*models.Departments, error)
	GetForFilterWithoutPagination(ctx *context.Context, req *dtos.GetFilter) ([]*models.Departments, error)
}

type Departments struct {
}

func NewDepartments() IDepartments {
	return &Departments{}
}

func (t *Departments) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".departments"
}

func (t *Departments) Upsert(ctx *context.Context, m *models.Departments) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *Departments) Get(ctx *context.Context, id string) (*models.Departments, error) {
	var result models.Departments
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *Departments) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.Departments{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *Departments) GetAll(ctx *context.Context) ([]*models.Departments, error) {
	var result []*models.Departments
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Departments) GetForFilterCount(ctx *context.Context, req *dtos.GetFilter) (int, error) {
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

func (t *Departments) GetForFilter(ctx *context.Context, req *dtos.GetFilter) ([]*models.Departments, error) {
	var result []*models.Departments
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

func (t *Departments) GetForFilterWithoutPagination(ctx *context.Context, req *dtos.GetFilter) ([]*models.Departments, error) {
	var result []*models.Departments
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

func (t *Departments) GetForRegionId(ctx *context.Context, regionId string) ([]*models.Departments, error) {
	var result []*models.Departments
	tx := ctx.DB.Table(t.getTable(ctx))
	tx = tx.Where("region_id =?", regionId)
	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}
	return result, nil
}
