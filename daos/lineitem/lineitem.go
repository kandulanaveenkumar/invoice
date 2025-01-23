package lineitem

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type ILineItem interface {
	Upsert(ctx *context.Context, m ...*models.LineItem) error
	Get(ctx *context.Context, id string) (*models.LineItem, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.LineItem, error)
	Delete(ctx *context.Context, id string) error
	DeleteByQuoteId(ctx *context.Context, quoteId string) error
	GetPartnersByQuoteId(ctx *context.Context, quoteId []string) ([]*models.LineItemPartners, error)
	GetLineItemsWithFilter(ctx *context.Context, filter *models.LiFields) ([]*models.LineItem, error)

	GetLineItemWithExchangeByRegionId(ctx *context.Context, lineItemId string, regionId string) (*models.LineItemWithExRate, error)
	GetLineItemsWithExchangeByQuoteId(ctx *context.Context, quoteId, regionId, lineItemId string, noRegionCheck bool) ([]*models.LineItemWithExRate, error)
	GetLineItemsWithExchangeForInvoice(ctx *context.Context, quoteId, regionId, customerId string, isMpbCustomer bool) ([]*models.LineItemWithExRate, error)

	UpdateInvoicedLineItems(ctx *context.Context, lineItemIds []string, isSellGenerated bool) error
	DeleteAll(ctx *context.Context, ids []string) error
	GetDashboardPartnerLineItems(ctx *context.Context, partneId string, regionId string) ([]*models.LineItem, error)
	GetDashboardCustomerLineItems(ctx *context.Context, quoteIds []string, regionId string) ([]*models.LineItem, error)
	ValidateLiForGeneratedInvoice(ctx *context.Context, lineItemIds []string) (bool, error)
	GetDetailsWithRefLineItemIds(ctx *context.Context, refLiIds []string) ([]*models.LineItem, error)

	GetLineItemsWithExchangeByQuoteIdAndLineitemId(ctx *context.Context, quoteId string, regionId string, lineItemIds []string, noRegionCheck bool) ([]*models.LineItemWithExRate, error)
}

type LineItem struct {
}

func NewLineItem() ILineItem {
	return &LineItem{}
}

func (t *LineItem) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "line_items"
}

func (t *LineItem) Upsert(ctx *context.Context, m ...*models.LineItem) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *LineItem) Get(ctx *context.Context, id string) (*models.LineItem, error) {
	var result models.LineItem
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *LineItem) Delete(ctx *context.Context, id string) error {
	var result models.LineItem
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete lineitem.", zap.Error(err))
		return err
	}

	return err
}

func (t *LineItem) DeleteByQuoteId(ctx *context.Context, quoteId string) error {
	var result models.LineItem
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "quote_id = ?", quoteId).Error
	if err != nil {
		ctx.Log.Error("Unable to delete lineitems.", zap.Error(err))
		return err
	}

	return err
}

func (t *LineItem) GetAll(ctx *context.Context, ids []string) ([]*models.LineItem, error) {
	var result []*models.LineItem
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get lineitems.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitems.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *LineItem) GetPartnersByQuoteId(ctx *context.Context, quoteIds []string) ([]*models.LineItemPartners, error) {
	var result []*models.LineItemPartners
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("DISTINCT quote_id, partner_id").Find(&result, "quote_id IN (?)", quoteIds).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (l *LineItem) GetLineItemWithExchangeByRegionId(ctx *context.Context, lineItemId string, regionId string) (*models.LineItemWithExRate, error) {
	var result *models.LineItemWithExRate
	query := ctx.DB.Table(l.getTable(ctx)).Select("line_items.*", "buy_ex.exchange_rate AS buy_exchange_rate", "sell_ex.exchange_rate as exchange_rate").
		Joins("JOIN "+constants.GetExchangeRatesTableForLineItem(l.getTable(ctx))+" buy_ex ON buy_ex.line_item_id = line_item.id AND buy_ex.type = 'buyrate' AND buy_ex.region_id = ?", regionId).
		Joins("JOIN "+constants.GetExchangeRatesTableForLineItem(l.getTable(ctx))+" sell_ex ON sell_ex.line_item_id = line_item.id AND sell_ex.type = 'sellrate' AND sell_ex.region_id = ?", regionId).
		Where("id = ? AND region_id = ?", lineItemId, regionId)

	err := query.First(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem for quote_id", zap.Error(err), zap.Any("qid", lineItemId))
		return nil, err
	}
	return result, nil
}

func (l *LineItem) GetLineItemsWithExchangeByQuoteId(ctx *context.Context, quoteId, regionId, lineItemId string, noRegionCheck bool) ([]*models.LineItemWithExRate, error) {
	var result []*models.LineItemWithExRate
	query := ctx.DB.Debug().Table(l.getTable(ctx)).Select("line_items.*", "buy_ex.exchange_rate AS buy_exchange_rate", "sell_ex.exchange_rate as exchange_rate").
		Joins("LEFT JOIN "+constants.GetExchangeRatesTableForLineItem(l.getTable(ctx))+" buy_ex ON buy_ex.line_item_id = line_items.id AND buy_ex.type = 'buyrate' AND buy_ex.region_id = ?", regionId).
		Joins("LEFT JOIN "+constants.GetExchangeRatesTableForLineItem(l.getTable(ctx))+" sell_ex ON sell_ex.line_item_id = line_items.id AND sell_ex.type = 'sellrate' AND sell_ex.region_id = ?", regionId).
		Where("quote_id = ?", quoteId)

	if !noRegionCheck {
		query.Where("line_items.region_id = ?", regionId)
	}

	if lineItemId != "" {
		query.Where("line_items.id = ?", lineItemId)
	}

	query.Group(" line_items.id,buy_ex.exchange_rate,sell_ex.exchange_rate ").Order("line_items.created_at")

	err := query.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem for quote_id", zap.Error(err), zap.Any("qid", quoteId))
		return nil, err
	}
	return result, nil
}

func (l *LineItem) GetLineItemsWithFilter(ctx *context.Context, filter *models.LiFields) ([]*models.LineItem, error) {
	var result []*models.LineItem
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(l.getTable(ctx)).Where(filter).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem for quote_id", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (l *LineItem) UpdateInvoicedLineItems(ctx *context.Context, lineItemIds []string, isSellGenerated bool) error {
	err := ctx.DB.Table(l.getTable(ctx)).Where("id IN ?", lineItemIds).Update("is_sell_invoice_generated", true).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem for quote_id", zap.Error(err))
		return err
	}
	return nil
}

func (l *LineItem) GetLineItemsWithExchangeForInvoice(ctx *context.Context, quoteId, regionId, customerId string, isMpbCustomer bool) ([]*models.LineItemWithExRate, error) {
	var result []*models.LineItemWithExRate
	query := ctx.DB.Debug().WithContext(ctx).Table(l.getTable(ctx)).Select("line_items.*", "buy_ex.exchange_rate AS buy_exchange_rate", "sell_ex.exchange_rate as exchange_rate").
		Joins("JOIN "+constants.GetExchangeRatesTableForLineItem(l.getTable(ctx))+" buy_ex ON buy_ex.line_item_id = line_items.id AND buy_ex.type = 'buyrate' AND buy_ex.region_id = ?", regionId).
		Joins("JOIN "+constants.GetExchangeRatesTableForLineItem(l.getTable(ctx))+" sell_ex ON sell_ex.line_item_id = line_items.id AND sell_ex.type = 'sellrate' AND sell_ex.region_id = ?", regionId).
		Where("quote_id = ? and line_items.region_id = ?", quoteId, regionId)

	if isMpbCustomer {
		query = query.Where("mpb_company_id = ?", customerId)
	} else {
		query = query.Where("line_items.is_mpb_addon = ?", false)
	}

	err := query.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem for quote_id", zap.Error(err), zap.Any("qid", quoteId))
		return nil, err
	}
	return result, nil
}

func (l *LineItem) DeleteAll(ctx *context.Context, ids []string) error {
	var result models.LineItemExchangeRates
	err := ctx.DB.Debug().Table(l.getTable(ctx)).Delete(&result, "id IN ?", ids).Error
	if err != nil {
		ctx.Log.Error("Unable to delete approvalpendinglineitemexchangerates.", zap.Error(err))
		return err
	}

	return err
}

func (t *LineItem) GetDashboardPartnerLineItems(ctx *context.Context, partneId string, regionId string) ([]*models.LineItem, error) {
	var result []*models.LineItem
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Joins("INNER JOIN shipments s ON line_items.quote_id = s.quote_id").
		Where("s.is_deleted = false AND line_items.partner_id = ? AND line_items.region_id = ?", partneId, regionId).
		Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem for quote_id", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (t *LineItem) GetDashboardCustomerLineItems(ctx *context.Context, quoteIds []string, regionId string) ([]*models.LineItem, error) {
	var result []*models.LineItem
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Joins("JOIN shipments s ON line_items.quote_id = s.quote_id").
		Where("s.is_deleted = false AND line_items.quote_id IN (?) and line_items.region_id = ?", quoteIds, regionId).
		Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem for quote_id", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (t *LineItem) ValidateLiForGeneratedInvoice(ctx *context.Context, lineItemIds []string) (bool, error) {
	var cnt int64
	err := ctx.DB.Table(t.getTable(ctx)).Raw(`
	SELECT COUNT(DISTINCT id) as count_ids
	FROM (
		SELECT DISTINCT id
		FROM "public"."line_items"
		WHERE id IN ? 
		AND is_sell_invoice_generated = ?
	) as filtered_ids
	HAVING COUNT(DISTINCT id) = ?
`, lineItemIds, true, len(lineItemIds)).Scan(&cnt).Error

	if err != nil {
		return false, err
	}
	return cnt > 0, err
}

func (t *LineItem) GetDetailsWithRefLineItemIds(ctx *context.Context, refLiIds []string) ([]*models.LineItem, error) {
	var result []*models.LineItem
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("id,is_sell_invoice_generated,ref_line_item_id").Where("ref_line_item_id IN ?", refLiIds).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (l *LineItem) GetLineItemsWithExchangeByQuoteIdAndLineitemId(ctx *context.Context, quoteId string, regionId string, lineItemIds []string, noRegionCheck bool) ([]*models.LineItemWithExRate, error) {
	var result []*models.LineItemWithExRate
	query := ctx.DB.Debug().Table(l.getTable(ctx)).Select("line_items.*", "buy_ex.exchange_rate AS buy_exchange_rate", "sell_ex.exchange_rate as exchange_rate").
		Joins("LEFT JOIN "+constants.GetExchangeRatesTableForLineItem(l.getTable(ctx))+" buy_ex ON buy_ex.line_item_id = line_items.id AND buy_ex.type = 'buyrate' AND buy_ex.region_id = ?", regionId).
		Joins("LEFT JOIN "+constants.GetExchangeRatesTableForLineItem(l.getTable(ctx))+" sell_ex ON sell_ex.line_item_id = line_items.id AND sell_ex.type = 'sellrate' AND sell_ex.region_id = ?", regionId).
		Where("quote_id = ?", quoteId)

	if !noRegionCheck {
		query.Where("line_items.region_id = ?", regionId)
	}

	if len(lineItemIds) > 0 {
		query.Where("line_items.id IN (?)", lineItemIds)
	}

	query.Group(" line_items.id,buy_ex.exchange_rate,sell_ex.exchange_rate ").Order("line_items.created_at")

	err := query.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem for quote_id", zap.Error(err), zap.Any("qid", quoteId))
		return nil, err
	}
	return result, nil
}
