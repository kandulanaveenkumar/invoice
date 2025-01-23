package billsregenerate

import (
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IBillsRegenerate interface {
	Upsert(ctx *context.Context, m ...*models.BillsRegenerate) error
	UpsertBillsGenerate(ctx *context.Context, masterBillsRegenerate *models.BillsRegenerate, shimentId uuid.UUID) error
	Get(ctx *context.Context, id string) (*models.BillsRegenerate, error)
	GetBillsRegenerate(ctx *context.Context, id uuid.UUID, billId uuid.UUID, query string) (*models.BillsRegenerate, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.BillsRegenerate, error)
	Delete(ctx *context.Context, id string) error
}

type BillsRegenerate struct {
}

func NewBillsRegenerate() IBillsRegenerate {
	return &BillsRegenerate{}
}

func (t *BillsRegenerate) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "bills_regenerate"
}

func (t *BillsRegenerate) Upsert(ctx *context.Context, m ...*models.BillsRegenerate) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *BillsRegenerate) UpsertBillsGenerate(ctx *context.Context, masterBillsRegenerate *models.BillsRegenerate, shimentId uuid.UUID) error {

	currentTime := time.Now().UTC()

	if masterBillsRegenerate.Id == uuid.Nil {
		masterBillsRegenerate.Id = uuid.New()
		masterBillsRegenerate.CreatedAt = currentTime
		masterBillsRegenerate.CreatedBy = ctx.Account.ID
	}

	masterBillsRegenerate.ShipmentId = shimentId
	masterBillsRegenerate.UpdatedAt = currentTime
	masterBillsRegenerate.UpdatedBy = ctx.Account.ID

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(masterBillsRegenerate).Error
	if err != nil {
		ctx.Log.Error("Failed to upsert BillsRegenerate", zap.Error(err))
		return err
	}

	return nil
}

func (t *BillsRegenerate) Get(ctx *context.Context, id string) (*models.BillsRegenerate, error) {
	var result models.BillsRegenerate
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get billsregenerate.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *BillsRegenerate) GetBillsRegenerate(ctx *context.Context, id uuid.UUID, billId uuid.UUID, query string) (*models.BillsRegenerate, error) {

	var result models.BillsRegenerate
	tx := ctx.DB.WithContext(ctx.Request.Context()).
		Debug().
		Table(t.getTable(ctx))

	if id != uuid.Nil {
		tx = tx.Where("id = ?", id)
	}
	if billId != uuid.Nil {
		tx = tx.Where("bill_id = ?", billId)
	}
	if query != "" {
		tx = tx.Where(query)
	}

	tx.Order("created_at DESC")
	err := tx.First(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get billsregenerate.", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

func (t *BillsRegenerate) Delete(ctx *context.Context, id string) error {
	var result models.BillsRegenerate
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete billsregenerate.", zap.Error(err))
		return err
	}

	return err
}

func (t *BillsRegenerate) GetAll(ctx *context.Context, ids []string) ([]*models.BillsRegenerate, error) {
	var result []*models.BillsRegenerate
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get billsregenerates.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get billsregenerates.", zap.Error(err))
		return nil, err
	}

	return result, err
}
