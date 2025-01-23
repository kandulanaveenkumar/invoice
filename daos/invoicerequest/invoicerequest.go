package invoicerequest

import (
	"errors"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IInvoiceRequest interface {
	Upsert(ctx *context.Context, m ...*models.InvoiceRequest) error
	Update(ctx *context.Context, m *models.InvoiceRequest) error
	Get(ctx *context.Context, id string) (*models.InvoiceRequest, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.InvoiceRequest, error)
	Delete(ctx *context.Context, id string) error
	ValidateAndAuditRequest(ctx *context.Context, shipmentId, regionId uuid.UUID, invType string) (*models.InvoiceRequest, error)
	UpdateStatus(ctx *context.Context, err error) error
	MarkCompleted(ctx *context.Context) error
}

type InvoiceRequest struct {
}

func NewInvoiceRequest() IInvoiceRequest {
	return &InvoiceRequest{}
}

func (t *InvoiceRequest) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "invoice_requests"
}

func (t *InvoiceRequest) Upsert(ctx *context.Context, m ...*models.InvoiceRequest) error {
	return ctx.DB.Table(t.getTable(ctx)).Save(m).Error
}

func (t *InvoiceRequest) Update(ctx *context.Context, m *models.InvoiceRequest) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Debug().Where("id = ?", m.ID).Updates(m).Error
}

func (t *InvoiceRequest) Get(ctx *context.Context, id string) (*models.InvoiceRequest, error) {
	var result models.InvoiceRequest
	err := ctx.DB.Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoicerequest.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *InvoiceRequest) Delete(ctx *context.Context, id string) error {
	var result models.InvoiceRequest
	err := ctx.DB.Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete invoicerequest.", zap.Error(err))
		return err
	}

	return err
}

func (t *InvoiceRequest) GetAll(ctx *context.Context, ids []string) ([]*models.InvoiceRequest, error) {
	var result []*models.InvoiceRequest
	if len(ids) == 0 {
		err := ctx.DB.Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get invoicerequests.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoicerequests.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *InvoiceRequest) ValidateAndAuditRequest(ctx *context.Context, shipmentId, regionId uuid.UUID, invType string) (*models.InvoiceRequest, error) {

	requestId, err := uuid.Parse(ctx.RefID)
	if err != nil {
		ctx.Log.Error("invalid request id", zap.Any("id", ctx.RefID), zap.Error(err))
	}

	tx := ctx.DB.Table(t.getTable(ctx)).
		Clauses(clause.Locking{Strength: "SHARE", Options: "NOWAIT"}).
		Begin()
	err = tx.Error
	if err != nil {
		ctx.Log.Error("unable to acquire table lock", zap.Error(err))
		return nil, err
	}

	var result *models.InvoiceRequest
	err = tx.Where("shipment_id = ? AND region_id = ? AND invoice_type = ? AND is_completed = false", shipmentId, regionId, invType).
		Scan(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		ctx.Log.Error("Unable to get invoicerequest.", zap.Error(err))
		tx.Rollback()
		return nil, err
	}
	if result != nil {
		err = errors.New("another request is already in progress")
		ctx.Log.Info("request already in progress for similar invoice", zap.Any("active_request_id", result.ID))
		tx.Rollback()
		return nil, err
	}

	err = tx.Save(&models.InvoiceRequest{
		ID:          requestId,
		CreatedBy:   ctx.Account.ID,
		ShipmentId:  shipmentId,
		RegionId:    regionId,
		InvoiceType: invType,
		IsCompleted: false,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to upsert invoicerequest.", zap.Error(err))
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return result, err
}

func (t *InvoiceRequest) UpdateStatus(ctx *context.Context, resErr error) error {
	columns := make(map[string]interface{})
	columns["is_successful"] = true
	if resErr != nil {
		columns["is_successful"] = false
		columns["error_message"] = resErr.Error()
	}

	err := ctx.DB.Table(t.getTable(ctx)).Where("id = ?", ctx.RefID).UpdateColumns(columns).Error
	if err != nil {
		ctx.Log.Error("unable to update invoicerequest", zap.Error(err))
	}
	return err
}

func (t *InvoiceRequest) MarkCompleted(ctx *context.Context) error {
	err := ctx.DB.Table(t.getTable(ctx)).Where("id = ?", ctx.RefID).Update("is_completed", true).Error
	if err != nil {
		ctx.Log.Error("unable to update invoicerequest", zap.Error(err))
	}
	return err
}
