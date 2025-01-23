package versionedlineitems

import (
	"database/sql"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IVersionedLineItems interface {
	Upsert(ctx *context.Context, m ...*models.VersionedLineItem) error
	Get(ctx *context.Context, id string) (*models.LineItem, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.LineItem, error)
	Delete(ctx *context.Context, id string) error
	DeleteByQuoteId(ctx *context.Context, quoteId string) error
	GetWIthExchangeRatesByVersionAndId(ctx *context.Context, id, regionId string, version int64) (*models.LineItemWithExRate, error)
	GetVersionedLineItemsWithExchangeRates(ctx *context.Context, quoteId, regionId string, version int64) ([]*models.LineItemWithExRate, error)
	GetMaxLineItemVersion(ctx *context.Context, quoteID string) (int64, error)
	GetAllMaxVersionedLiWitExRates(ctx *context.Context, ids []string, regionID string) ([]*models.LineItemWithExRate, error)
	GetMaxLineItemsForTimeline(ctx *context.Context, qid string, buyRegionId string, version int64) ([]*models.TimeLineLineItems, error)
}

type VersionedLineItems struct {
}

func NewVersionedLineItems() IVersionedLineItems {
	return &VersionedLineItems{}
}

func (t *VersionedLineItems) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "versioned_line_items"
}

func (t *VersionedLineItems) Upsert(ctx *context.Context, m ...*models.VersionedLineItem) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *VersionedLineItems) Get(ctx *context.Context, id string) (*models.LineItem, error) {
	var result models.LineItem
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get versioned_line_items.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *VersionedLineItems) Delete(ctx *context.Context, id string) error {
	var result models.LineItem
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete versioned_line_item.", zap.Error(err))
		return err
	}

	return err
}

func (t *VersionedLineItems) DeleteByQuoteId(ctx *context.Context, quoteId string) error {
	var result models.LineItem
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "quote_id = ?", quoteId).Error
	if err != nil {
		ctx.Log.Error("Unable to delete versioned_line_items.", zap.Error(err))
		return err
	}

	return err
}

func (t *VersionedLineItems) GetAll(ctx *context.Context, ids []string) ([]*models.LineItem, error) {
	var result []*models.LineItem
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get versioned_line_itemss.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get versioned_line_itemss.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *VersionedLineItems) GetWIthExchangeRatesByVersionAndId(ctx *context.Context, id, regionId string, version int64) (*models.LineItemWithExRate, error) {
	var result *models.LineItemWithExRate
	err := ctx.DB.Table(t.getTable(ctx)).Select("line_items.*", "buy_ex.exchange_rate AS buy_exchange_rate", "sell_ex.exchange_rate as exchange_rate").
		Joins("JOIN "+constants.GetExchangeRatesTableForLineItem(t.getTable(ctx))+" buy_ex ON buy_ex.line_item_id = line_item.id AND buy_ex.version = line_item.version AND buy_ex.type = 'buyrate' AND buy_ex.region_id = ?", regionId).
		Joins("JOIN "+constants.GetExchangeRatesTableForLineItem(t.getTable(ctx))+" sell_ex ON sell_ex.line_item_id = line_item.id AND sell_ex.version = line_item.version AND sell_ex.type = 'sellrate' AND sell_ex.region_id = ?", regionId).
		First(&result, "id = ? AND version = ?", id, version).Error
	if err != nil {
		ctx.Log.Error("Unable to get versioned_line_items.", zap.Error(err))
		return nil, err
	}

	return result, err
}
func (t *VersionedLineItems) GetVersionedLineItemsWithExchangeRates(ctx *context.Context, quoteId, regionId string, version int64) ([]*models.LineItemWithExRate, error) {
	var result []*models.LineItemWithExRate

	q := ctx.DB.Table(t.getTable(ctx)).Select("versioned_line_items.*", "buy_ex.exchange_rate AS buy_exchange_rate", "sell_ex.exchange_rate as exchange_rate").
		Joins("JOIN "+constants.GetVersionedExchangeRatesTableForLineItem(t.getTable(ctx))+" buy_ex ON buy_ex.line_item_id = versioned_line_items.id AND buy_ex.version = versioned_line_items.version AND buy_ex.type = 'buyrate' AND buy_ex.region_id = ?", regionId).
		Joins("JOIN "+constants.GetVersionedExchangeRatesTableForLineItem(t.getTable(ctx))+" sell_ex ON sell_ex.line_item_id = versioned_line_items.id AND sell_ex.version = versioned_line_items.version AND sell_ex.type = 'sellrate' AND sell_ex.region_id = ?", regionId).
		Where("versioned_line_items.region_id = ?", regionId)

	if quoteId != "" {
		q = q.Where("quote_id = ?", quoteId)
	}
	if version != 0 {
		q = q.Where("versioned_line_items.version = ?", version)
	}

	err := q.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get versioned_line_items.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *VersionedLineItems) GetMaxLineItemVersion(ctx *context.Context, quoteID string) (int64, error) {

	var maxVersion sql.NullInt64
	err := ctx.DB.Table(t.getTable(ctx)).Select("max(version)").Where("quote_id = ?", quoteID).Scan(&maxVersion).Error
	if err != nil {
		ctx.Log.Error("error fetching max version", zap.Error(err))
		return 0, err
	}
	return maxVersion.Int64, nil
}

func (t *VersionedLineItems) GetAllMaxVersionedLiWitExRates(ctx *context.Context, ids []string, regionId string) ([]*models.LineItemWithExRate, error) {
	var result []*models.LineItemWithExRate

	var maxVersions []models.MaxVersion
	q := ctx.DB.Table(t.getTable(ctx)).Select("id, MAX(version) as max_version")

	if len(ids) != 0 {
		q = q.Where("id IN ?", ids)
	}

	q = q.Group("id").Find(&maxVersions)
	if q.Error != nil {
		ctx.Log.Error("Unable to get versioned_line_items.", zap.Error(q.Error))
		return nil, q.Error
	}

	err := ctx.DB.Table(t.getTable(ctx)).Debug().Select("versioned_line_items.*", "buy_ex.exchange_rate AS buy_exchange_rate", "sell_ex.exchange_rate as exchange_rate").
		Joins("JOIN (?) AS mv ON versioned_line_items.id = mv.id AND versioned_line_items.version = mv.max_version", q).
		Joins("JOIN "+constants.GetVersionedExchangeRatesTableForLineItem(t.getTable(ctx))+" buy_ex ON buy_ex.line_item_id = versioned_line_items.id AND buy_ex.version = versioned_line_items.version AND buy_ex .region_id = ? AND buy_ex.type = 'buyrate'", regionId).
		Joins("JOIN "+constants.GetVersionedExchangeRatesTableForLineItem(t.getTable(ctx))+" sell_ex ON sell_ex.line_item_id = versioned_line_items.id AND sell_ex.version = versioned_line_items.version AND sell_ex .region_id = ? AND sell_ex.type = 'sellrate'", regionId).
		Find(&result).Error

	if err != nil {
		ctx.Log.Error("Unable to get versioned_line_items.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *VersionedLineItems) GetMaxLineItemsForTimeline(ctx *context.Context, qid string, buyRegionId string, version int64) ([]*models.TimeLineLineItems, error) {
	var liVersions []*models.TimeLineLineItems

	err := ctx.DB.Debug().Table(t.getTable(ctx)).Select(
		"id as line_item_id",
		"versioned_line_items.version",
		"versioned_line_items.sell * sell_ex.exchange_rate * versioned_line_items.units + ((versioned_line_items.tax / 100) * versioned_line_items.sell * sell_ex.exchange_rate * versioned_line_items.units) AS total_sell",
		"versioned_line_items.buy * buy_ex.exchange_rate * versioned_line_items.units + ((versioned_line_items.buy_tax / 100) * versioned_line_items.buy * buy_ex.exchange_rate * versioned_line_items.units) AS total_buy",
	).Joins(
		"INNER JOIN "+constants.GetVersionedExchangeRatesTableForLineItem(t.getTable(ctx))+" buy_ex ON buy_ex.line_item_id = versioned_line_items.id AND buy_ex.version = versioned_line_items.version AND buy_ex.region_id = ? AND buy_ex.type = 'buyrate'", buyRegionId,
	).Joins(
		"INNER JOIN "+constants.GetVersionedExchangeRatesTableForLineItem(t.getTable(ctx))+" sell_ex ON sell_ex.line_item_id = versioned_line_items.id AND sell_ex.version = versioned_line_items.version AND sell_ex.region_id = ? AND sell_ex.type = 'sellrate'", buyRegionId,
	).Where(
		"versioned_line_items.quote_id = ? AND versioned_line_items.region_id = ? AND versioned_line_items.version = ? AND versioned_line_items.sub_type != 'Tax' AND versioned_line_items.sub_type != 'Handling Charge'", qid, buyRegionId, version,
	).Find(&liVersions).Error

	if err != nil {
		ctx.Log.Error("Unable to get versioned_line_items.", zap.Error(err))
		return nil, err
	}

	return liVersions, nil
}
