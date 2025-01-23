package invoice

import (
	"strings"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IInvoice interface {
	Upsert(ctx *context.Context, m ...*models.Invoice) error
	UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.Invoice) error
	Get(ctx *context.Context, id string) (*models.Invoice, error)
	GetAll(ctx *context.Context, ids []string, offset, limit int) ([]*models.Invoice, error)
	Delete(ctx *context.Context, id string) error

	GetWithFilter(ctx *context.Context, filter *models.Invoice) ([]*models.Invoice, error)
	GetInvoiceAmount(ctx *context.Context, req models.InvoiceAmount) ([]*models.Invoice, error)
	CheckForGeneratedInvoice(ctx *context.Context, shipmentId, regionId string, invoiceTypes []string) (bool, error)
	GetForInvoiceByCompanyIds(ctx *context.Context, shipmentId string, companyIds []string, invoiceType string, regionId string) ([]*models.Invoice, error)
	GetByInvoiceNumber(ctx *context.Context, number string) (*models.Invoice, error)
	GetByLineItemIds(ctx *context.Context, lineItemIds []string, invoiceType string) ([]*models.Invoice, error)
	GetDebitNoteItems(ctx *context.Context, shipmentId string) (*models.CreditDebitNoteLineItems, error)
	GetCreditNoteItems(ctx *context.Context, shipmentId string) (*models.CreditDebitNoteLineItems, error)
	GetTotalDebitNoteItems(ctx *context.Context, lineItemIds []string, regionId string) (*models.TotalCreditDebitAmount, error)
	GetTotalCreditNoteItems(ctx *context.Context, lineItemIds []string, regionId string) (*models.TotalCreditDebitAmount, error)
	Update(ctx *context.Context, m *models.Invoice) error
	GetTotalCount(ctx *context.Context) (int, error)
	GetIcaVendorInvoice(ctx *context.Context, shipmentId string, voucherId string) (*models.Invoice, error)
}

type Invoice struct {
}

func NewInvoice() IInvoice {
	return &Invoice{}
}

func (t *Invoice) getTable(ctx *context.Context) string {
	if ctx.TenantID == "" {
		ctx.TenantID = "public"
	}
	return ctx.TenantID + "." + "invoices"
}

func (t *Invoice) getInvoiceLineItemsTable(ctx *context.Context) string {
	if ctx.TenantID == "" {
		ctx.TenantID = "public"
	}
	return ctx.TenantID + "." + "invoice_line_items"
}

func (t *Invoice) Upsert(ctx *context.Context, m ...*models.Invoice) error {

	err := ctx.DB.Table(t.getTable(ctx)).Save(m).Error
	if err != nil {
		ctx.Log.Error("unable to upsert invoice", zap.Error(err))
	}
	return err
}

func (t *Invoice) UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.Invoice) error {

	err := tx.Table(t.getTable(ctx)).Save(m).Error
	if err != nil {
		ctx.Log.Error("unable to upsert invoice with tx", zap.Error(err))
	}
	return err
}

func (t *Invoice) Get(ctx *context.Context, id string) (*models.Invoice, error) {

	var result models.Invoice

	err := ctx.DB.Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoice.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *Invoice) Delete(ctx *context.Context, id string) error {

	var result models.Invoice

	err := ctx.DB.Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete invoice.", zap.Error(err))
		return err
	}

	return err
}

func (t *Invoice) GetAll(ctx *context.Context, ids []string, offset, limit int) ([]*models.Invoice, error) {
	var result []*models.Invoice

	if len(ids) == 0 {

		tx := ctx.DB.Table(t.getTable(ctx)).Where("vat_treatment = ''").Offset(offset)

		if limit != 0 {
			tx = tx.Limit(limit)
		}

		err := tx.Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get invoices.", zap.Error(err))
			return nil, err
		}
		return result, err
	}

	err := ctx.DB.Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoices.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Invoice) GetInvoiceAmount(ctx *context.Context, req models.InvoiceAmount) ([]*models.Invoice, error) {

	invoices := []*models.Invoice{}

	tx := ctx.DB.Debug().Table(t.getTable(ctx))
	if len(req.VoucherTypes) > 0 {
		tx.Where("voucher_type in (?)", req.VoucherTypes)
	}

	if len(req.InvoiceTypes) > 0 {
		tx.Where("invoice_type in (?)", req.InvoiceTypes)
	}

	if len(req.ShipmentIds) > 0 {
		tx.Where("shipment_id in (?)", req.ShipmentIds)
	}

	if len(req.NotVoucherTypes) > 0 {
		tx.Where("voucher_type NOT in (?)", req.NotVoucherTypes)
	}

	err := tx.First(&invoices).Error
	if err != nil {
		return nil, err
	}
	return invoices, nil
}

func (t *Invoice) GetWithFilter(ctx *context.Context, filter *models.Invoice) ([]*models.Invoice, error) {
	var result []*models.Invoice

	tx := ctx.DB.Debug().Table(t.getTable(ctx))

	if filter != nil && filter.CompanyId != uuid.Nil {
		tx.Where("company_id = ?", filter.CompanyId)
	}

	if filter != nil && filter.ShipmentID != uuid.Nil {
		tx.Where("shipment_id = ?", filter.ShipmentID)
	}

	if filter != nil && filter.InvoiceType != "" {
		tx.Where("invoice_type = ?", filter.InvoiceType)
	}

	if filter != nil && filter.RegionId != uuid.Nil {
		tx.Where("region_id = ?", filter.RegionId)
	}

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoices.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Invoice) CheckForGeneratedInvoice(ctx *context.Context, shipmentId, regionId string, invoiceTypes []string) (bool, error) {
	tx := ctx.DB.Table(t.getTable(ctx)).Where("shipment_id = ?", shipmentId)

	if regionId != "" {
		tx.Where("region_id = ?", regionId)
	}

	if len(invoiceTypes) > 0 {
		tx.Where("invoice_type IN (?)", invoiceTypes)
	}

	var isGenerated bool
	err := ctx.DB.Raw("SELECT EXISTS (?)", tx).Scan(&isGenerated).Error
	if err != nil {
		return false, err
	}
	return isGenerated, err
}

func (t *Invoice) GetForInvoiceByCompanyIds(ctx *context.Context, shipmentId string, companyIds []string, invoiceType string, regionId string) ([]*models.Invoice, error) {
	var result []*models.Invoice

	tx := ctx.DB.Debug().Table(t.getTable(ctx))

	if len(companyIds) > 0 {
		tx.Where("company_id IN (?)", companyIds)
	}

	if shipmentId != "" {
		tx.Where("shipment_id =?", shipmentId)
	}

	if invoiceType != "" {
		tx.Where("invoice_type =?", invoiceType)
	}

	if regionId != "" {
		tx.Where("region_id =?", regionId)
	}

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoices.", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (t *Invoice) GetByInvoiceNumber(ctx *context.Context, number string) (*models.Invoice, error) {

	var result models.Invoice

	err := ctx.DB.Table(t.getTable(ctx)).First(&result, "no = ?", number).Error
	if err != nil {
		ctx.Log.Error("unable to get invoice details by invoice number", zap.Error(err), zap.Any("invoice number", number))
		return nil, err
	}

	return &result, err
}

func (t *Invoice) GetByLineItemIds(ctx *context.Context, lineItemIds []string, invoiceType string) ([]*models.Invoice, error) {
	var invoices []*models.Invoice
	var invoiceLineItems []*models.InvoiceLineItem

	err := ctx.DB.Debug().Table(t.getTable(ctx)+" i").
		Select("i.id, i.voucher_id, i.invoiced_date, i.due_on, i.no, i.partner_inv_docs").
		Joins("JOIN "+t.getInvoiceLineItemsTable(ctx)+" ili ON i.id = ili.invoice_id").
		Where("ili.line_item_id IN (?) AND i.invoice_type = ?", lineItemIds, invoiceType).
		Group("i.id").
		Order("i.created_at").
		Scan(&invoices).Error
	if err != nil {
		ctx.Log.Error("error fetching invoices", zap.Error(err))
		return nil, err
	}

	invoiceIds := make([]string, len(invoices))
	for i, inv := range invoices {
		invoiceIds[i] = inv.ID.String()
	}

	err = ctx.DB.Debug().Table(t.getInvoiceLineItemsTable(ctx)).
		Select("id, invoice_id, line_item_id").
		Where("line_item_id IN (?) AND invoice_id IN (?)", lineItemIds, invoiceIds).
		Order("created_at").
		Scan(&invoiceLineItems).Error
	if err != nil {
		ctx.Log.Error("error fetching invoice line items", zap.Error(err))
		return nil, err
	}

	invoiceMap := make(map[uuid.UUID]*models.Invoice)
	for _, invoice := range invoices {
		invoice.LineItems = []*models.InvoiceLineItem{}
		invoiceMap[invoice.ID] = invoice
	}

	for _, lineItem := range invoiceLineItems {
		if invoice, exists := invoiceMap[lineItem.InvoiceId]; exists {
			invoice.LineItems = append(invoice.LineItems, lineItem)
		}
	}

	return invoices, nil
}

func (t *Invoice) GetDebitNoteItems(ctx *context.Context, shipmentId string) (*models.CreditDebitNoteLineItems, error) {
	var debitItems []*models.DebitNoteLineItem
	result := ctx.DB.Debug().Raw(`
        SELECT il.rate AS debit, il.line_item_id AS id, il.quantity, il.exchange_rate,il.currency
        FROM invoices i
        JOIN invoice_line_items il ON i.id = il.invoice_id
        WHERE i.shipment_id = ? AND i.invoice_type = 'debit_note'
    `, shipmentId).Scan(&debitItems)

	if result.Error != nil {
		return nil, result.Error
	}

	creditDebitItems := &models.CreditDebitNoteLineItems{
		DebitNoteLineItems: debitItems,
	}

	return creditDebitItems, nil
}

func (t *Invoice) GetCreditNoteItems(ctx *context.Context, shipmentId string) (*models.CreditDebitNoteLineItems, error) {
	var creditItems []*models.CreditNoteLineItem
	result := ctx.DB.Debug().Raw(`
        SELECT il.rate AS credit, il.line_item_id AS id, il.quantity, il.exchange_rate, il.currency
        FROM invoices i
        JOIN invoice_line_items il ON i.id = il.invoice_id
        WHERE i.shipment_id = ? AND i.invoice_type = 'credit_note'
    `, shipmentId).Scan(&creditItems)

	if result.Error != nil {
		return nil, result.Error
	}

	creditDebitItems := &models.CreditDebitNoteLineItems{
		CreditNoteLineItems: creditItems,
	}

	return creditDebitItems, nil
}

func (t *Invoice) GetTotalCreditNoteItems(ctx *context.Context, lineItemIds []string, regionId string) (*models.TotalCreditDebitAmount, error) {
	combinedlineItemIds := "( '" + strings.Join(lineItemIds, "','") + "') "
	var Items *models.TotalCreditDebitAmount

	result := ctx.DB.Debug().Raw(`
	SELECT 
    coalesce(SUM(il.rate * il.quantity * lier.exchange_rate)) AS total_amount,
    coalesce(SUM(il.rate* il.quantity * lier.exchange_rate * il.tax_percentage/100)) AS total_rate
    FROM 
    invoices i
    JOIN 
    invoice_line_items il 
    ON i.id = il.invoice_id
	JOIN
    line_item_exchange_rates lier 
    on lier.line_item_id  = il.line_item_id
    WHERE 
    il.line_item_id IN ` + combinedlineItemIds + `
	AND i.invoice_type = 'credit_note' and   lier.region_id ='` + regionId + `'
    AND lier."type" ='sellrate'
	`).Scan(&Items)

	if result.Error != nil {
		return nil, result.Error
	}

	return Items, nil
}
func (t *Invoice) GetTotalDebitNoteItems(ctx *context.Context, lineItemIds []string, regionId string) (*models.TotalCreditDebitAmount, error) {
	combinedlineItemIds := "( '" + strings.Join(lineItemIds, "','") + "') "
	var Items *models.TotalCreditDebitAmount
	result := ctx.DB.Debug().Raw(`
	SELECT 
    coalesce(SUM(il.rate * il.quantity * lier.exchange_rate)) AS total_amount,
    coalesce(SUM(il.rate* il.quantity * lier.exchange_rate * il.tax_percentage/100)) AS total_rate
    FROM 
    invoices i
    JOIN 
    invoice_line_items il 
    ON i.id = il.invoice_id
	JOIN
    line_item_exchange_rates lier 
    on lier.line_item_id  = il.line_item_id
    WHERE 
    il.line_item_id IN ` + combinedlineItemIds + `
	AND i.invoice_type = 'debit_note' and   lier.region_id ='` + regionId + `'
    AND lier."type" ='buyrate'
	`).Scan(&Items)

	if result.Error != nil {
		return nil, result.Error
	}

	return Items, nil
}

func (t *Invoice) Update(ctx *context.Context, m *models.Invoice) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Debug().Where("id = ?", m.ID).Updates(m).Error
}

func (t Invoice) GetTotalCount(ctx *context.Context) (int, error) {

	var count int64
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (t *Invoice) GetIcaVendorInvoice(ctx *context.Context, shipmentId string, voucherId string) (*models.Invoice, error) {

	var result models.Invoice

	err := ctx.DB.Table(t.getTable(ctx)).First(&result, "shipment_id = ? AND voucher_id = ?", shipmentId, voucherId).Error
	if err != nil {
		ctx.Log.Error("unable to get invoice details by shipment and voucher ID", zap.Error(err), zap.Any("voucher id", voucherId))
		return nil, err
	}

	return &result, err
}
