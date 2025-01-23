package invoicesequence

import (
	"errors"
	"strings"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IInvoiceSequence interface {
	Upsert(ctx *context.Context, m ...*models.InvoiceSequence) error
	Get(ctx *context.Context, id string) (*models.InvoiceSequence, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.InvoiceSequence, error)
	Delete(ctx *context.Context, id string) error
	NewInvoiceNumber(ctx *context.Context, invoiceType, voucherType string, regionId uuid.UUID) (int64, error)
}

type InvoiceSequence struct {
}

func NewInvoiceSequence() IInvoiceSequence {
	return &InvoiceSequence{}
}

func (t *InvoiceSequence) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "invoice_sequences"
}

func (t *InvoiceSequence) Upsert(ctx *context.Context, m ...*models.InvoiceSequence) error {
	return ctx.DB.Table(t.getTable(ctx)).Save(m).Error
}

func (t *InvoiceSequence) Get(ctx *context.Context, id string) (*models.InvoiceSequence, error) {
	var result models.InvoiceSequence
	err := ctx.DB.Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoicesequence.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *InvoiceSequence) Delete(ctx *context.Context, id string) error {
	var result models.InvoiceSequence
	err := ctx.DB.Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete invoicesequence.", zap.Error(err))
		return err
	}

	return err
}

func (t *InvoiceSequence) GetAll(ctx *context.Context, ids []string) ([]*models.InvoiceSequence, error) {
	var result []*models.InvoiceSequence
	if len(ids) == 0 {
		err := ctx.DB.Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get invoicesequences.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get invoicesequences.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (i *InvoiceSequence) NewInvoiceNumber(ctx *context.Context, invoiceType string, voucherType string, regionId uuid.UUID) (int64, error) {
	seqName, err := getSequenceName(ctx, invoiceType, voucherType)
	if err != nil {
		return 0, err
	}

	sequence := models.InvoiceSequence{
		RegionId: regionId,
		Name:     seqName,
	}

	// Fetch the existing sequence
	err = ctx.DB.Table(i.getTable(ctx)).
		Where("region_id = ? AND name = ?", regionId, seqName).
		First(&sequence).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// If not found, set initial values for a new sequence
		sequence.No = 1
		sequence.CreatedBy = ctx.Account.ID
		sequence.UpdatedBy = ctx.Account.ID
	} else if err != nil {
		ctx.Log.Error("Failed to retrieve invoice sequence.", zap.Error(err))
		return 0, err
	} else {
		// If found, increment the sequence number and update fields
		sequence.No++
		sequence.UpdatedBy = ctx.Account.ID
	}

	// Save the record
	if err := ctx.DB.Save(&sequence).Error; err != nil {
		ctx.Log.Error("Failed to save invoice sequence.", zap.Error(err))
		return 0, err
	}

	return sequence.No, nil
}

func getSequenceName(ctx *context.Context, invoiceType, voucherType string) (string, error) {
	seqName := ""
	if strings.Contains(invoiceType, constants.CustomerProforma) {
		seqName = "est_"
	}
	switch voucherType {
	case "INV":
		seqName += constants.InvoiceNoCustomerInvoice
	case "VI":
		seqName += constants.InvoiceNoVendorInvoice
	case "BOS":
		seqName += constants.InvoiceNoBillOfSupply
	case "RN":
		seqName += constants.InvoiceNoReimbursementNote
	case "CN":
		seqName += constants.InvoiceNoCreditNote
	case "DN":
		seqName += constants.InvoiceNoDebitNote
	case "OSCN":
		seqName += constants.InvoiceNoCreditNote
	default:
		ctx.Log.Error("unknown types for sequence name", zap.Any("invoice_type", invoiceType), zap.Any("voucher_type", voucherType))
		return "", errors.New("unknown types for sequence name")
	}
	return seqName, nil
}
