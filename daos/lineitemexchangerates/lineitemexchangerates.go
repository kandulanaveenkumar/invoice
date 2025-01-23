package lineitemexchangerates

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type ILineItemExchangeRates interface {
	Upsert(ctx *context.Context, m ...*models.LineItemExchangeRates) error
	Get(ctx *context.Context, id string) ([]*models.LineItemExchangeRates, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.LineItemExchangeRates, error)
	Delete(ctx *context.Context, id string) error
	GetAllLineItemsExchangeRateWithType(ctx *context.Context, lineItemIds []string, regionId, exRatetype string) ([]*models.LineItemExchangeRates, error)
	DeleteAll(ctx *context.Context, ids []string) error
	GetBulkLineItemsExchangeRatesWithFilter(ctx *context.Context, lis []string, filter *models.LineItemExchangeRates) ([]*models.LineItemExchangeRates, error)
}

type LineItemExchangeRates struct {
}

func NewLineItemExchangeRates() ILineItemExchangeRates {
	return &LineItemExchangeRates{}
}

func (t *LineItemExchangeRates) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "line_item_exchange_rates"
}

func (t *LineItemExchangeRates) Upsert(ctx *context.Context, m ...*models.LineItemExchangeRates) error {
	constraint := "line_item_exchange_rates_pkey"
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Clauses(clause.OnConflict{OnConstraint: constraint, UpdateAll: true}).Save(m).Error
}

func (t *LineItemExchangeRates) Get(ctx *context.Context, id string) ([]*models.LineItemExchangeRates, error) {
	var result []*models.LineItemExchangeRates
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result, "line_item_id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitemexchangerates.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *LineItemExchangeRates) Delete(ctx *context.Context, id string) error {
	var result models.LineItemExchangeRates
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "line_item_id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete lineitemexchangerates.", zap.Error(err))
		return err
	}

	return err
}

func (t *LineItemExchangeRates) GetAll(ctx *context.Context, ids []string) ([]*models.LineItemExchangeRates, error) {
	var result []*models.LineItemExchangeRates
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get lineitemexchangeratess.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("line_item_id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitemexchangeratess.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *LineItemExchangeRates) GetAllLineItemsExchangeRate(ctx *context.Context, lineItemIds []string) ([]*models.LineItemExchangeRates, error) {
	var result []*models.LineItemExchangeRates
	if len(lineItemIds) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get lineitemexchangeratess.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("line_item_id IN ?", lineItemIds).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitemexchangeratess.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *LineItemExchangeRates) GetAllLineItemsExchangeRateWithType(ctx *context.Context, lineItemIds []string, regionId, exRatetype string) ([]*models.LineItemExchangeRates, error) {
	var result []*models.LineItemExchangeRates
	err := ctx.DB.Table(t.getTable(ctx)).Where("line_item_id IN ?", lineItemIds).Where("type = ? AND region_id = ?", exRatetype, regionId).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitemexchangeratess.", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (t *LineItemExchangeRates) DeleteAll(ctx *context.Context, ids []string) error {
	var result models.LineItemExchangeRates
	err := ctx.DB.Debug().Table(t.getTable(ctx)).Delete(&result, "line_item_id IN ?", ids).Error
	if err != nil {
		ctx.Log.Error("Unable to delete lineitemexchangeratess.", zap.Error(err))
		return err
	}

	return err
}

func (t *LineItemExchangeRates) GetBulkLineItemsExchangeRatesWithFilter(ctx *context.Context, lis []string, filter *models.LineItemExchangeRates) ([]*models.LineItemExchangeRates, error) {
	var result []*models.LineItemExchangeRates
	q := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Where("line_item_id IN (?)", lis)
	if filter.Currency != "" {
		q = q.Where("currency = ?", filter.Currency)
	}
	if filter.Type != "" {
		q = q.Where("type = ?", filter.Type)
	}
	if filter.RegionId != uuid.Nil && filter.RegionId.String() != "" {
		q = q.Where("region_id = ?", filter.RegionId)
	}
	err := q.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
