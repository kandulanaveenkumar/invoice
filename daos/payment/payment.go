package payment

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IPayment interface {
	Upsert(ctx *context.Context, m ...*models.Payment) error
	Get(ctx *context.Context, id string) (*models.Payment, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.Payment, error)
	Delete(ctx *context.Context, id string) error
}

type Payment struct {
}

func NewPayment() IPayment {
	return &Payment{}
}

func (t *Payment) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "payments"
}

func (t *Payment) Upsert(ctx *context.Context, m ...*models.Payment) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *Payment) Get(ctx *context.Context, id string) (*models.Payment, error) {
	var result models.Payment
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get payment.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *Payment) Delete(ctx *context.Context, id string) error {
	var result models.Payment
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete payment.", zap.Error(err))
		return err
	}

	return err
}

func (t *Payment) GetAll(ctx *context.Context, ids []string) ([]*models.Payment, error) {
	var result []*models.Payment
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get payments.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get payments.", zap.Error(err))
		return nil, err
	}

	return result, err
}
