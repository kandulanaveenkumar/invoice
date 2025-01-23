package invoicelineitem

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IInvoiceLineItem interface {
	Upsert(ctx *context.Context, m ...*models.InvoiceLineItem) error
	UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.InvoiceLineItem) error
	Get(ctx *context.Context, id string) (*models.InvoiceLineItem, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.InvoiceLineItem, error)
	Delete(ctx *context.Context, id string) error
	GetInvoiceWithVoucherTypes(ctx *context.Context, lineItemId uuid.UUID, voucherType []string, invoiceType string, partnerId uuid.UUID, shipmentId uuid.UUID, query string) (*models.Invoice, error)
	GetLineItemsCountWithShipmentId(ctx *context.Context, invoiceType string, shipmentId uuid.UUID, query string) (*models.DistinctLineitemCount, error)
	GetGeneratedLineItems(ctx *context.Context, req models.InvoiceAmount) ([]*models.InvoiceLineItem, error)
	GetLineItemForLatestInvoice(ctx *context.Context, lineItemId uuid.UUID, voucherType []string, invoiceType, billToAccountId, baseCurrency, exchangeRate, query string) (*models.InvoiceLineItem, error)
	GetLineItemsForInvoice(ctx *context.Context, lineItemId uuid.UUID, voucherType []string, invoiceType, billToAccountId, shipmentId, status, query string) ([]*models.InvoiceLineItem, error)
	CheckForGeneratedInvoice(ctx *context.Context, lineItemIds []string, invoiceType string, single bool) (interface{}, error)
	GetInvoicedAmounts(ctx *context.Context, lineItemIds []string, invoiceType string) ([]*models.LineItemInvoicedAmountWithType, error)
	GetForInvoiceId(ctx *context.Context, invoiceId string) ([]*models.InvoiceLineItem, error)
	GetTotalSoFar(ctx *context.Context, lineItemId uuid.UUID, invoiceType string) ([]*models.InvoiceLineItem, error)
	GetInvoiceLineItemsFilter(ctx *context.Context, filters *models.InvoiceLineItemsFilters) ([]*models.InvoiceLineItem, error)
	DeleteInvoiceLineItemsByInvoiceId(ctx *context.Context, InvoiceId string) error
}

type InvoiceLineItem struct {
}

func NewInvoiceLineItem() IInvoiceLineItem {
	return &InvoiceLineItem{}
}

func (t *InvoiceLineItem) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "invoice_line_items"
}

func (t *InvoiceLineItem) Upsert(ctx *context.Context, m ...*models.InvoiceLineItem) error {
	err := ctx.DB.Table(t.getTable(ctx)).Save(m).Error
	if err != nil {
		ctx.Log.Error("unable to upsert invoicelineitems", zap.Error(err))
	}
	return err
}

func (t *InvoiceLineItem) UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.InvoiceLineItem) error {
	err := tx.Debug().Table(t.getTable(ctx)).Save(m).Error
	if err != nil {
		ctx.Log.Error("unable to upsert invoicelineitems with tx", zap.Error(err))
	}
	return err
}

func (t *InvoiceLineItem) Get(ctx *context.Context, id string) (*models.InvoiceLineItem, error) {
	var result models.InvoiceLineItem
	err := ctx.DB.Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoicelineitem.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *InvoiceLineItem) Delete(ctx *context.Context, id string) error {
	var result models.InvoiceLineItem
	err := ctx.DB.Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete invoicelineitem.", zap.Error(err))
		return err
	}

	return err
}

func (t *InvoiceLineItem) GetAll(ctx *context.Context, ids []string) ([]*models.InvoiceLineItem, error) {
	var result []*models.InvoiceLineItem
	if len(ids) == 0 {
		err := ctx.DB.Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get invoicelineitems.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoicelineitems.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *InvoiceLineItem) GetInvoiceWithVoucherTypes(ctx *context.Context, lineItemId uuid.UUID, voucherType []string, invoiceType string, partnerId uuid.UUID, shipmentId uuid.UUID, query string) (*models.Invoice, error) {

	invoices := models.Invoice{}

	tx := ctx.DB.Table(t.getTable(ctx)).Select("invoices.*").Joins("JOIN invoice_line_items ON invoices.id = invoice_line_items.invoice_id")

	if lineItemId != uuid.Nil {
		tx.Where("invoice_line_items.line_item_id = ?", lineItemId)
	}
	if len(voucherType) > 0 {
		tx.Where("invoices.voucher_type in (?) ", voucherType)
	}
	if invoiceType != "" {
		tx.Where("invoices.invoice_type = ?", invoiceType)
	}
	if partnerId != uuid.Nil {
		tx.Where("invoice_line_items.partner_id = ?", partnerId)
	}
	if shipmentId != uuid.Nil {
		tx.Where("invoices.booking_id = ?", shipmentId)
	}
	if query != "" {
		tx.Where(query)
	}

	tx.Order("invoice_line_items.created_at DESC")
	err := tx.First(&invoices).Error
	if err != nil {
		return nil, err
	}

	return &invoices, nil
}

func (t *InvoiceLineItem) GetLineItemsCountWithShipmentId(ctx *context.Context, invoiceType string, shipmentId uuid.UUID, query string) (*models.DistinctLineitemCount, error) {

	total := &models.DistinctLineitemCount{}

	tx := ctx.DB.Table(t.getTable(ctx)).Select("DISTINCT count(invoice_line_items.line_item_id) as total").
		Joins("JOIN invoices ON invoices.id = invoice_line_items.invoice_id").
		Where("invoices.shipment_id = ?", shipmentId)

	if invoiceType != "" {
		tx.Where("invoices.invoice_type = ?", invoiceType)
	}

	if query != "" {
		tx.Where(query)
	}

	err := tx.Find(&total).Error
	if err != nil {
		return nil, err
	}

	return total, nil
}

func (t *InvoiceLineItem) GetLineItemsForInvoice(ctx *context.Context, lineItemId uuid.UUID, voucherType []string, invoiceType string, billToAccountId string, shipmentId, status, query string) ([]*models.InvoiceLineItem, error) {

	InvoicelineItem := []*models.InvoiceLineItem{}

	tx := ctx.DB.Table(t.getTable(ctx)).Joins("JOIN invoices ON invoices.id = invoice_line_items.invoice_id")

	if lineItemId != uuid.Nil {
		tx.Where("invoice_line_items.line_item_id = ?", lineItemId)
	}
	if len(voucherType) > 0 {
		tx.Where("invoices.voucher_type in (?)", voucherType)
	}
	if invoiceType != "" {
		tx.Where("invoices.invoice_type = ?", invoiceType)
	}
	if billToAccountId != "" {
		tx.Where("invoices.bill_to_account_id = ?", billToAccountId)
	}
	if shipmentId != "" {
		tx.Where("invoices.shipment_id = ?", shipmentId)
	}
	if status != "" {
		tx.Where("invoices.status = ?", status)

	}
	if query != "" {
		tx.Where(query)
	}

	tx.Order("created_at DESC")
	err := tx.Debug().Find(&InvoicelineItem).Error
	if err != nil {
		return nil, err
	}

	return InvoicelineItem, nil
}

func (t *InvoiceLineItem) GetLineItemForLatestInvoice(ctx *context.Context, lineItemId uuid.UUID, voucherType []string, invoiceType string, billToAccountId string, baseCurrency, exchangeRate, query string) (*models.InvoiceLineItem, error) {

	lineItem := models.InvoiceLineItem{}

	tx := ctx.DB.Table(t.getTable(ctx)).Select("invoice_line_items.*").Joins("JOIN invoices ON invoices.id = invoice_line_items.invoice_id")

	if lineItemId != uuid.Nil {
		tx.Where("invoice_line_items.line_item_id = ?", lineItemId)
	}
	if len(voucherType) > 0 {
		tx.Where("invoices.voucher_type in (?)", voucherType)
	}
	if invoiceType != "" {
		tx.Where("invoices.invoice_type = ?", invoiceType)
	}

	if billToAccountId != "" {
		tx.Where("invoices.bill_to_account_id = ?", billToAccountId)
	}

	if baseCurrency != "" {
		tx.Where("invoices.base_currency = ?", baseCurrency)
	}

	if exchangeRate != "" {
		tx.Where("invoice_line_items.exchange_rate = ?", exchangeRate)
	}

	if query != "" {
		tx.Where(query)
	}

	tx.Order("invoice_line_items.created_at DESC")
	err := tx.First(&lineItem).Error
	if err != nil {
		return nil, err
	}

	return &lineItem, nil
}

func (t *InvoiceLineItem) GetGeneratedLineItems(ctx *context.Context, req models.InvoiceAmount) ([]*models.InvoiceLineItem, error) {
	invoiceLineItems := []*models.InvoiceLineItem{}
	tx := ctx.DB.Table(t.getTable(ctx)).Select("invoice_line_items.*,invoices.booking_id,invoices.booking_id,invoices.number").Joins("JOIN invoices ON invoices.id = invoice_line_items.invoice_id")
	if len(req.ShipmentIds) > 0 {
		tx.Where("invoices.shipment_id::text in (?)", req.ShipmentIds)
	}
	if len(req.InvoiceTypes) > 0 {
		tx.Where("invoices.invoice_type in (?)", req.InvoiceTypes)
	}
	if len(req.VoucherTypes) > 0 {
		tx.Where("invoices.voucher_type in (?)", req.VoucherTypes)
	}
	if len(req.NotVoucherTypes) > 0 {
		tx.Where("invoices.voucher_type NOT in (?)", req.NotVoucherTypes)
	}
	err := tx.Find(&invoiceLineItems).Error
	if err != nil {
		return nil, err
	}
	return invoiceLineItems, nil
}

// CheckForGeneratedInvoice checks if invoices have been generated either for all line items (bool) or individually for each line item (map).
// If `single` is true, it returns a single boolean indicating whether any of the line items have a generated invoice.
// If `single` is false, it returns a map with line_item_id as the key and a boolean indicating if the invoice has been generated for that line item.
func (t *InvoiceLineItem) CheckForGeneratedInvoice(ctx *context.Context, lineItemIds []string, invoiceType string, single bool) (interface{}, error) {
	if single {
		// Single result: Check if any invoice has been generated for the line items
		tx := ctx.DB.Table(t.getTable(ctx)).Joins("JOIN invoices ON invoices.id = invoice_line_items.invoice_id").
			Where("line_item_id IN (?) AND invoice_type = ?", lineItemIds, invoiceType)

		var isGenerated bool
		err := ctx.DB.Raw("SELECT EXISTS (?)", tx).Scan(&isGenerated).Error
		if err != nil {
			return false, err
		}
		return isGenerated, nil
	}

	// Multiple results: Check individually for each line item
	var results []struct {
		LineItemId  string
		IsGenerated bool
	}

	// Fixed query with MAX to handle multiple entries and booleans
	tx := ctx.DB.Table(t.getTable(ctx)).
		Select("line_item_id, MAX(CASE WHEN EXISTS (SELECT 1 FROM invoices WHERE invoices.id = invoice_line_items.invoice_id AND invoice_type = ?) THEN 1 ELSE 0 END) as is_generated", invoiceType).
		Where("line_item_id IN (?)", lineItemIds).
		Group("line_item_id")

	// Execute the query and store results
	err := tx.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Populate the map with lineItemId to generation status
	isGeneratedMap := make(map[string]bool)
	for _, res := range results {
		isGeneratedMap[res.LineItemId] = res.IsGenerated
	}

	return isGeneratedMap, nil
}

func (t *InvoiceLineItem) GetInvoicedAmounts(ctx *context.Context, lineItemIds []string, invoiceType string) ([]*models.LineItemInvoicedAmountWithType, error) {
	var invoicedAmounts []*models.LineItemInvoicedAmountWithType
	var invType []string

	if invoiceType == "both" {
		invType = append(invType, constants.CreditNote, constants.DebitNote)
	} else {
		invType = append(invType, invoiceType)
	}

	err := ctx.DB.Table(t.getTable(ctx)).Select("invoice_line_items.line_item_id", "invoice_line_items.invoice_id", "invoice_line_items.currency", "invoice_line_items.tax_amount as amount", "invoices.invoice_type").Joins("JOIN invoices ON invoices.id = invoice_line_items.invoice_id").
		Where("line_item_id IN (?) AND invoices.invoice_type IN (?)", lineItemIds, invType).Find(&invoicedAmounts).Error
	if err != nil {
		return nil, err
	}
	return invoicedAmounts, err
}

func (t *InvoiceLineItem) GetForInvoiceId(ctx *context.Context, invoiceId string) ([]*models.InvoiceLineItem, error) {

	invoiceLineItems := []*models.InvoiceLineItem{}

	tx := ctx.DB.Debug().Table(t.getTable(ctx))
	if invoiceId != "" {
		tx.Where("invoice_id = ?", invoiceId)
	}

	err := tx.Find(&invoiceLineItems).Error
	if err != nil {
		return nil, err
	}

	return invoiceLineItems, nil

}

func (t *InvoiceLineItem) GetTotalSoFar(ctx *context.Context, lineItemId uuid.UUID, invoiceType string) ([]*models.InvoiceLineItem, error) {

	invoicelineItems := []*models.InvoiceLineItem{}

	tx := ctx.DB.Table(t.getTable(ctx)).Joins("JOIN invoices ON invoices.id = invoice_line_items.invoice_id").Where("invoice_line_items.line_item_id = ? AND invoices.invoice_type = ?", lineItemId, invoiceType)

	err := tx.Find(&invoicelineItems).Error
	if err != nil {
		return nil, err
	}

	return invoicelineItems, nil
}

func (t *InvoiceLineItem) GetInvoiceLineItemsFilter(ctx *context.Context, filters *models.InvoiceLineItemsFilters) ([]*models.InvoiceLineItem, error) {
	invoiceLineItems := []*models.InvoiceLineItem{}
	tx := ctx.DB.Table(t.getTable(ctx)).Debug().Joins("JOIN invoices ON invoices.id = invoice_line_items.invoice_id")

	if len(filters.InvoiceIds) > 0 {
		tx.Where("invoice_line_items.invoice_id in (?)", filters.InvoiceIds)
	}

	if filters.InvoiceType != "" {
		tx.Where("invoices.invoice_type = ?", filters.InvoiceType)
	}

	if len(filters.LineItemIds) > 0 {
		tx.Where("invoice_line_items.line_item_id in (?)", filters.LineItemIds)
	}

	err := tx.Find(&invoiceLineItems).Error
	if err != nil {
		return nil, err
	}
	return invoiceLineItems, nil

}

func (t *InvoiceLineItem) DeleteInvoiceLineItemsByInvoiceId(ctx *context.Context, InvoiceId string) error {
	var result models.InvoiceLineItem
	err := ctx.DB.Table(t.getTable(ctx)).Delete(&result, "invoice_id = ?", InvoiceId).Error
	if err != nil {
		ctx.Log.Error("Unable to delete invoicelineitem.", zap.Error(err))
		return err
	}

	return err
}
