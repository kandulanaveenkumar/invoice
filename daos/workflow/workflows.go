package workflow

import (
	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IWorkflows interface {
	Upsert(ctx *context.Context, m *models.Workflows) error
	Get(ctx *context.Context, id string) (*models.Workflows, error)
	GetAll(ctx *context.Context) ([]*models.Workflows, error)
	Delete(ctx *context.Context, id string) error
	GetForFilterCount(ctx *context.Context, req *dtos.GetFilter) (int, error)
	GetForFilter(ctx *context.Context, req *dtos.GetFilter) ([]*models.Workflows, error)
	GetForShipment(ctx *context.Context, condition *dtos.WorkflowShipmentConditions) ([]*models.Workflows, error)
	GetFor(ctx *context.Context, condition *dtos.WorkflowConditions) ([]*models.Workflows, error)
	GetForFilterWithoutPagination(ctx *context.Context, req *dtos.GetFilter) ([]*models.Workflows, error)
}

type Workflows struct {
}

func NewWorkflows() IWorkflows {
	return &Workflows{}
}

func (t *Workflows) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".workflows"
}

func (t *Workflows) Upsert(ctx *context.Context, m *models.Workflows) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *Workflows) Get(ctx *context.Context, id string) (*models.Workflows, error) {
	var result models.Workflows
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *Workflows) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.Workflows{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *Workflows) GetFor(ctx *context.Context, condition *dtos.WorkflowConditions) ([]*models.Workflows, error) {
	var result []*models.Workflows
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if condition.Service != "" {
		tx = tx.Where(" ? =  ANY(services) OR cardinality(services) = 0", condition.Service)
	}

	if condition.IncoTerm != "" {
		tx = tx.Where(" ? =  ANY(inco_terms) OR cardinality(inco_terms) = 0", condition.IncoTerm)
	}

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Workflows) GetForShipment(ctx *context.Context, condition *dtos.WorkflowShipmentConditions) ([]*models.Workflows, error) {
	var result []*models.Workflows
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if condition.IncoTerm != "" {
		tx = tx.Where(" ? =  ANY(inco_terms) OR inco_terms = '{Select-All}'", condition.IncoTerm)
	}

	if condition.HsCode != "" {
		tx = tx.Where(" ? =  ANY(hs_codes) OR hs_codes = '{Select-All}'", condition.HsCode)
	}

	if condition.RegionId != "" {
		tx = tx.Where("region_id = ?", condition.RegionId)
	}

	if condition.Pol != "" {
		tx = tx.Where(" ? =  ANY(pols) OR pols = '{Select-All}'", condition.Pol)
	}

	if condition.Pod != "" {
		tx = tx.Where(" ? =  ANY(pods) OR pods = '{Select-All}'", condition.Pod)
	}

	err := tx.Where("enabled = true").Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Workflows) GetAll(ctx *context.Context) ([]*models.Workflows, error) {
	var result []*models.Workflows
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Workflows) GetForFilterCount(ctx *context.Context, req *dtos.GetFilter) (int, error) {
	var cnt int
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("count(*)")
	if req.Q != "" {
		tx = tx.Where("name ILIKE ?", req.Q+"%")
	}

	if req.ID != "" {
		tx = tx.Where("id = ?", req.ID)
	}

	if req.RegionID != "" {
		tx = tx.Where("region_id = ?", req.ID)
	}

	err := tx.Scan(&cnt).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return 0, err
	}

	return cnt, err
}

func (t *Workflows) GetForFilter(ctx *context.Context, req *dtos.GetFilter) ([]*models.Workflows, error) {
	var result []*models.Workflows
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if req.Q != "" {
		tx = tx.Where("name ILIKE ?", req.Q+"%")
	}

	if req.ID != "" {
		tx = tx.Where("id = ?", req.ID)
	}

	if req.RegionID != "" {
		tx = tx.Where("region_id = ?", req.ID)
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

func (t *Workflows) GetForFilterWithoutPagination(ctx *context.Context, req *dtos.GetFilter) ([]*models.Workflows, error) {
	var result []*models.Workflows
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
