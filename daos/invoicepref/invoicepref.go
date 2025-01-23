package invoicepref

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type IInvoicePref interface {
	Upsert(ctx *context.Context, m ...*models.InvoicePref) error
	Update(ctx *context.Context, m *models.InvoicePref) error
	Get(ctx *context.Context, shipmentId, regionId, companyId, prefType string) (*models.InvoicePref, error)
	GetAll(ctx *context.Context, shipmentId, prefType string) ([]*models.InvoicePref, error)
	Delete(ctx *context.Context, id string) error
}

type InvoicePref struct {
}

func NewInvoicePref() IInvoicePref {
	return &InvoicePref{}
}

func (t *InvoicePref) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "invoice_prefs"
}

func (t *InvoicePref) Upsert(ctx *context.Context, m ...*models.InvoicePref) error {
	for _, invoicePref := range m {
		if err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
			Clauses(clause.OnConflict{
				UpdateAll: true,
				Columns:   []clause.Column{{Name: "region_id"}, {Name: "shipment_id"}, {Name: "type"}, {Name: "company_id"}},
			}).
			Create(invoicePref).Error; err != nil {
			return err
		}
	}
	return nil
}

func (t *InvoicePref) Update(ctx *context.Context, m *models.InvoicePref) error {
	return ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Updates(m).Error
}

func (t *InvoicePref) Get(ctx *context.Context, shipmentId, regionId, companyId, prefType string) (*models.InvoicePref, error) {
	var result models.InvoicePref
	err := ctx.DB.Table(t.getTable(ctx)).First(&result, "shipment_id = ? and region_id = ? and type = ? and company_id = ?", shipmentId, regionId, prefType, companyId).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoicepref.", zap.Error(err))
	}

	return &result, err
}

func (t *InvoicePref) Delete(ctx *context.Context, id string) error {
	var result models.InvoicePref
	err := ctx.DB.Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete invoicepref.", zap.Error(err))
		return err
	}

	return err
}

func (t *InvoicePref) GetAll(ctx *context.Context, shipmentId, prefType string) ([]*models.InvoicePref, error) {
	var result []*models.InvoicePref

	tx := ctx.DB.Table(t.getTable(ctx))

	if shipmentId != "" {
		tx.Where("shipment_id = ?", shipmentId)
	}

	if prefType != "" {
		tx.Where("type = ?", prefType)
	}

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoiceprefs.", zap.Error(err))
		return nil, err
	}

	return result, err
}
