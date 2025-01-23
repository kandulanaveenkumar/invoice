package approvalpendinglineitems

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IApprovalPendingLineItems interface {
	Upsert(ctx *context.Context, m ...*models.ApprovalPendingLineItem) error
	Get(ctx *context.Context, id string) (*models.ApprovalPendingLineItem, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.ApprovalPendingLineItem, error)
	Delete(ctx *context.Context, id string) error
	GetApprovalPendingLisWithExRatesByQuoteId(ctx *context.Context, quoteId, regionId string) ([]*models.ApprovalPendingLiWithExRate, error)
	GetApprovalPendingLisByQuoteId(ctx *context.Context, quoteId string, fresh *bool) ([]*models.ApprovalPendingLineItem, error)
	DeleteAll(ctx *context.Context, ids []string) error
	GetApprovalPendingLisByQuoteIdAndInvoiceNumber(ctx *context.Context, quoteId, invoiceNumber string) ([]*models.ApprovalPendingLineItem, error)
	GetApprovalPendingLisWithExchangeByQuoteId(ctx *context.Context, quoteId, regionId, lineItemId string, noRegionCheck bool) ([]*models.ApprovalPendingLineItemWithExRate, error)
	GetApprovalPendingLisByLineItemId(ctx *context.Context, lineItemId string) (*models.ApprovalPendingLineItem, error)
}

type ApprovalPendingLineItems struct {
}

func NewApprovalPendingLineItems() IApprovalPendingLineItems {
	return &ApprovalPendingLineItems{}
}

func (t *ApprovalPendingLineItems) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "approval_pending_line_items"
}

func (t *ApprovalPendingLineItems) Upsert(ctx *context.Context, m ...*models.ApprovalPendingLineItem) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *ApprovalPendingLineItems) Get(ctx *context.Context, id string) (*models.ApprovalPendingLineItem, error) {
	var result models.ApprovalPendingLineItem
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get approval_pending_line_items.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *ApprovalPendingLineItems) Delete(ctx *context.Context, id string) error {
	var result models.ApprovalPendingLineItem
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete approval_pending_line_items.", zap.Error(err))
		return err
	}

	return err
}

func (t *ApprovalPendingLineItems) GetAll(ctx *context.Context, ids []string) ([]*models.ApprovalPendingLineItem, error) {
	var result []*models.ApprovalPendingLineItem
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get approval_pending_line_itemss.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get approval_pending_line_itemss.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *ApprovalPendingLineItems) GetApprovalPendingLisWithExRatesByQuoteId(ctx *context.Context, quoteId, regionId string) ([]*models.ApprovalPendingLiWithExRate, error) {
	var result []*models.ApprovalPendingLiWithExRate
	query := ctx.DB.Debug().Table(t.getTable(ctx)).Select("approval_pending_line_items.*", "buy_ex.exchange_rate AS buy_exchange_rate", "sell_ex.exchange_rate as exchange_rate").
		Joins("LEFT JOIN "+constants.GetApprovalPendingExchangeRatesTableForLineItem(t.getTable(ctx))+" buy_ex ON buy_ex.line_item_id = approval_pending_line_items.line_item_id AND buy_ex.type = 'buyrate' AND buy_ex.region_id = ? AND buy_ex.is_fresh", regionId).
		Joins("LEFT JOIN "+constants.GetApprovalPendingExchangeRatesTableForLineItem(t.getTable(ctx))+" sell_ex ON sell_ex.line_item_id = approval_pending_line_items.line_item_id AND sell_ex.type = 'sellrate' AND sell_ex.region_id = ? AND sell_ex.is_fresh", regionId).
		Where("quote_id = ? and approval_pending_line_items.region_id = ?", quoteId, regionId).Where("approval_pending_line_items.is_fresh").Group(" approval_pending_line_items.id,approval_pending_line_items.line_item_id,buy_ex.exchange_rate,sell_ex.exchange_rate ").
		Order("approval_pending_line_items.created_at")

	err := query.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get lineitem for quote_id", zap.Error(err), zap.Any("qid", quoteId))
		return nil, err
	}
	return result, nil
}

func (t *ApprovalPendingLineItems) GetApprovalPendingLisByQuoteId(ctx *context.Context, quoteId string, fresh *bool) ([]*models.ApprovalPendingLineItem, error) {
	var result []*models.ApprovalPendingLineItem

	q := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("quote_id = ?", quoteId)

	if fresh != nil {
		q.Where("is_fresh = ?", fresh)
	}

	err := q.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get approval_pending_line_itemss.", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (t *ApprovalPendingLineItems) DeleteAll(ctx *context.Context, ids []string) error {
	var result models.LineItem
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id IN ?", ids).Error
	if err != nil {
		ctx.Log.Error("Unable to delete lineitem.", zap.Error(err))
		return err
	}

	return err
}

func (apli *ApprovalPendingLineItems) GetApprovalPendingLisByQuoteIdAndInvoiceNumber(ctx *context.Context, quoteId, invoiceNumber string) ([]*models.ApprovalPendingLineItem, error) {
	var result []*models.ApprovalPendingLineItem

	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(apli.getTable(ctx))

	if quoteId != "" {
		tx.Where("quote_id = ?", quoteId)
	}

	if invoiceNumber != "" {
		tx.Where("invoice_number = ?", invoiceNumber)
	}

	tx.Order("created_at desc")

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get approval_pending_line_itemss.", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (apli *ApprovalPendingLineItems) GetApprovalPendingLisWithExchangeByQuoteId(ctx *context.Context, quoteId, regionId, lineItemId string, noRegionCheck bool) ([]*models.ApprovalPendingLineItemWithExRate, error) {

	var result []*models.ApprovalPendingLineItemWithExRate
	query := ctx.DB.Debug().Table(apli.getTable(ctx)).Select("approval_pending_line_items.*", "buy_ex.exchange_rate AS buy_exchange_rate", "sell_ex.exchange_rate as exchange_rate").
		Joins("LEFT JOIN "+constants.GetApprovalPendingExchangeRatesTableForLineItem(apli.getTable(ctx))+" buy_ex ON buy_ex.line_item_id = approval_pending_line_items.line_item_id AND buy_ex.type = 'buyrate' AND buy_ex.region_id = ?", regionId).
		Joins("LEFT JOIN "+constants.GetApprovalPendingExchangeRatesTableForLineItem(apli.getTable(ctx))+" sell_ex ON sell_ex.line_item_id = approval_pending_line_items.line_item_id AND sell_ex.type = 'sellrate' AND sell_ex.region_id = ?", regionId).
		Where("quote_id = ?", quoteId)

	if !noRegionCheck {
		query.Where("approval_pending_line_items.region_id = ?", regionId)
	}

	if lineItemId != "" {
		query.Where("approval_pending_line_items.id = ?", lineItemId)
	}

	query.Group(" approval_pending_line_items.id,buy_ex.exchange_rate,sell_ex.exchange_rate ").Order("approval_pending_line_items.created_at")

	err := query.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get approval pending lineitems for quote_id", zap.Error(err), zap.Any("qid", quoteId))
		return nil, err
	}
	return result, nil
}

func (apli *ApprovalPendingLineItems) GetApprovalPendingLisByLineItemId(ctx *context.Context, lineItemId string) (*models.ApprovalPendingLineItem, error) {

	var result models.ApprovalPendingLineItem
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(apli.getTable(ctx)).First(&result, "line_item_id = ?", lineItemId).Error
	if err != nil {
		ctx.Log.Error("unable to get approval_pending_line_items", zap.Error(err), zap.Any("line item id", lineItemId))
		return nil, err
	}

	return &result, err
}
