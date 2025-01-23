package approvalpendinglineitemexchangerates

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type Iapprovalpendinglineitemexchangerates interface {
	Upsert(ctx *context.Context, m ...*models.ApprovalPendingLiExRates) error
	Get(ctx *context.Context, id string) (*models.ApprovalPendingLiExRates, error)
	GetAll(ctx *context.Context, ids []string, isFresh *bool) ([]*models.ApprovalPendingLiExRates, error)
	Delete(ctx *context.Context, id string) error
	DeleteAll(ctx *context.Context, ids []string) error
}

type approvalpendinglineitemexchangerates struct {
}

func Newapprovalpendinglineitemexchangerates() Iapprovalpendinglineitemexchangerates {
	return &approvalpendinglineitemexchangerates{}
}

func (t *approvalpendinglineitemexchangerates) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "approval_pending_line_item_exchange_rates"
}

func (t *approvalpendinglineitemexchangerates) Upsert(ctx *context.Context, m ...*models.ApprovalPendingLiExRates) error {
	return ctx.DB.Table(t.getTable(ctx)).Save(m).Error
}

func (t *approvalpendinglineitemexchangerates) Get(ctx *context.Context, id string) (*models.ApprovalPendingLiExRates, error) {
	var result models.ApprovalPendingLiExRates
	err := ctx.DB.Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get approvalpendinglineitemexchangerates.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *approvalpendinglineitemexchangerates) Delete(ctx *context.Context, id string) error {
	var result models.ApprovalPendingLiExRates
	err := ctx.DB.Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete approvalpendinglineitemexchangerates.", zap.Error(err))
		return err
	}

	return err
}

func (t *approvalpendinglineitemexchangerates) GetAll(ctx *context.Context, ids []string, isFresh *bool) ([]*models.ApprovalPendingLiExRates, error) {
	var result []*models.ApprovalPendingLiExRates
	if len(ids) == 0 {
		err := ctx.DB.Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get approvalpendinglineitemexchangeratess.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	q := ctx.DB.Table(t.getTable(ctx)).Where("line_item_id IN ?", ids)
	if isFresh != nil {
		q = q.Where("is_fresh = ?", isFresh)
	}
	err := q.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get approvalpendinglineitemexchangeratess.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *approvalpendinglineitemexchangerates) DeleteAll(ctx *context.Context, ids []string) error {
	var result models.ApprovalPendingLiExRates
	err := ctx.DB.Table(t.getTable(ctx)).Delete(&result, "line_item_id IN ?", ids).Error
	if err != nil {
		ctx.Log.Error("Unable to delete approvalpendinglineitemexchangerates.", zap.Error(err))
		return err
	}

	return err
}
