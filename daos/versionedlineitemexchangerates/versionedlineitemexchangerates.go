package versionedlineitemexchangerates

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type Iversionedlineitemexchangerates interface {
	Upsert(ctx *context.Context, m ...*models.VersionedLineItemExchangeRates) error
	Get(ctx *context.Context, id string) (*models.VersionedLineItemExchangeRates, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.VersionedLineItemExchangeRates, error)
	Delete(ctx *context.Context, id string) error
	DeleteAll(ctx *context.Context, ids []string) error
	GetVersionedExRatesWithIDs(ctx *context.Context, liIds []string, filter *models.VersionedLineItemExchangeRates) ([]*models.VersionedLineItemExchangeRates, error)
}

type versionedlineitemexchangerates struct {
}

func Newversionedlineitemexchangerates() Iversionedlineitemexchangerates {
	return &versionedlineitemexchangerates{}
}

func (t *versionedlineitemexchangerates) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "versioned_line_item_exchange_rates"
}

func (t *versionedlineitemexchangerates) Upsert(ctx *context.Context, m ...*models.VersionedLineItemExchangeRates) error {
	return ctx.DB.Table(t.getTable(ctx)).Save(m).Error
}

func (t *versionedlineitemexchangerates) Get(ctx *context.Context, id string) (*models.VersionedLineItemExchangeRates, error) {
	var result models.VersionedLineItemExchangeRates
	err := ctx.DB.Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get versionedlineitemexchangerates.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *versionedlineitemexchangerates) Delete(ctx *context.Context, id string) error {
	var result models.VersionedLineItemExchangeRates
	err := ctx.DB.Table(t.getTable(ctx)).Delete(&result, "line_item_id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete versionedlineitemexchangerates.", zap.Error(err))
		return err
	}

	return err
}

func (t *versionedlineitemexchangerates) DeleteAll(ctx *context.Context, ids []string) error {
	var result models.VersionedLineItemExchangeRates
	err := ctx.DB.Debug().Table(t.getTable(ctx)).Delete(&result, "line_item_id IN ?", ids).Error
	if err != nil {
		ctx.Log.Error("Unable to delete versionedlineitemexchangerates.", zap.Error(err))
		return err
	}

	return err
}

func (t *versionedlineitemexchangerates) GetAll(ctx *context.Context, ids []string) ([]*models.VersionedLineItemExchangeRates, error) {
	var result []*models.VersionedLineItemExchangeRates
	if len(ids) == 0 {
		err := ctx.DB.Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get versionedlineitemexchangeratess.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get versionedlineitemexchangeratess.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *versionedlineitemexchangerates) GetVersionedExRatesWithIDs(ctx *context.Context, liIds []string, filter *models.VersionedLineItemExchangeRates) ([]*models.VersionedLineItemExchangeRates, error) {
	var result []*models.VersionedLineItemExchangeRates
	q := ctx.DB.Table(t.getTable(ctx)).Where("line_item_id IN ?", liIds)
	if filter.Currency != "" {
		q = q.Where("currency = ?", filter.Currency)
	}
	if filter.Version != 0 {
		q = q.Where("version = ?", filter.Version)
	}
	err := q.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get versionedlineitemexchangeratess.", zap.Error(err))
		return nil, err
	}

	return result, err
}
