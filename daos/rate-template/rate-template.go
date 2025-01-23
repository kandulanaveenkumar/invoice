package ratetemplate

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IRateTemplate interface {
	Upsert(ctx *context.Context, m ...*models.RateTemplate) error
	Get(ctx *context.Context, id, rid string) (*models.RateTemplate, error)
	GetAll(ctx *context.Context, ids []string, service string, regionId string, companyId string) ([]*models.RateTemplate, error)
	Delete(ctx *context.Context, id string) error
}

type RateTemplate struct {
}

func NewRateTemplate() IRateTemplate {
	return &RateTemplate{}
}

func (t *RateTemplate) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "rate_templates"
}

func (t *RateTemplate) Upsert(ctx *context.Context, m ...*models.RateTemplate) error {
	return ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *RateTemplate) Get(ctx *context.Context, id, regionID string) (*models.RateTemplate, error) {
	var result models.RateTemplate

	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if id != "" {
		tx.Where("id = ?", id)
	}

	if regionID != "" {
		tx.Where("region_id = ?", regionID)
	}

	err := tx.Take(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rate template", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *RateTemplate) GetAll(ctx *context.Context, ids []string, shipmentType string, regionId string, companyId string) ([]*models.RateTemplate, error) {
	var result []*models.RateTemplate

	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if len(ids) > 0 {
		tx.Where("id IN ?", ids)
	}

	if shipmentType != "" {
		tx.Where("shipment_type", shipmentType)
	}

	if regionId != "" {
		tx.Where("region_id", regionId)
	}

	if companyId != "" {
		tx.Where("(array_length(applicability, 1) IS NULL OR ? = ANY(applicability))", companyId)
	}

	tx.Order("created_at")

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rate templates", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *RateTemplate) Delete(ctx *context.Context, id string) error {
	var result models.RateTemplate

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete rate template", zap.Error(err))
		return err
	}

	return err
}
