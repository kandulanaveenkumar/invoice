package rfqregionrestrictions

import (
	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type Irfqregionrestrictions interface {
	Upsert(ctx *context.Context, m ...*models.RfqRegionRestrictions) error
	Get(ctx *context.Context, id string) (*models.RfqRegionRestrictions, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.RfqRegionRestrictions, error)
	Delete(ctx *context.Context, id string) error
	GetByFilters(ctx *context.Context, filters *dtos.RfqRegionRestrictions) ([]*models.RfqRegionRestrictions, error)
}

type rfqregionrestrictions struct {
}

func Newrfqregionrestrictions() Irfqregionrestrictions {
	return &rfqregionrestrictions{}
}

func (t *rfqregionrestrictions) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "rfq_region_restrictions"
}

func (t *rfqregionrestrictions) Upsert(ctx *context.Context, m ...*models.RfqRegionRestrictions) error {
	return ctx.DB.Table(t.getTable(ctx)).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "region_id"}, {Name: "rfq_id"}, {Name: "quote_id"}},
			UpdateAll: true,
		}).Save(m).Error
}

func (t *rfqregionrestrictions) Get(ctx *context.Context, id string) (*models.RfqRegionRestrictions, error) {
	var result models.RfqRegionRestrictions
	err := ctx.DB.Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqregionrestrictions.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *rfqregionrestrictions) Delete(ctx *context.Context, id string) error {
	var result models.RfqRegionRestrictions
	err := ctx.DB.Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete rfqregionrestrictions.", zap.Error(err))
		return err
	}

	return err
}

func (t *rfqregionrestrictions) GetAll(ctx *context.Context, ids []string) ([]*models.RfqRegionRestrictions, error) {
	var result []*models.RfqRegionRestrictions
	if len(ids) == 0 {
		err := ctx.DB.Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get rfqregionrestrictionss.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqregionrestrictionss.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *rfqregionrestrictions) GetByFilters(ctx *context.Context, filters *dtos.RfqRegionRestrictions) ([]*models.RfqRegionRestrictions, error) {

	var result []*models.RfqRegionRestrictions
	q := ctx.DB.Table(t.getTable(ctx))

	if filters.QuoteID != "" {
		q = q.Where("quote_id = ?", filters.QuoteID)
	}

	if filters.RegionID != "" {
		q = q.Where("region_id = ?", filters.RegionID)
	}

	err := q.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqregionrestrictionss.", zap.Error(err))
		return nil, err
	}

	return result, err
}
