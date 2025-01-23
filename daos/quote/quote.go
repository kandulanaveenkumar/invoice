package quote

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IQuote interface {
	Upsert(ctx *context.Context, m ...*models.Quote) error
	Get(ctx *context.Context, id string) (*models.Quote, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.Quote, error)
	Delete(ctx *context.Context, id string) error
	Update(ctx *context.Context, m *models.Quote) error
}

type Quote struct {
}

func NewQuote() IQuote {
	return &Quote{}
}

func (t *Quote) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "quotes"
}

func (t *Quote) Upsert(ctx *context.Context, m ...*models.Quote) error {
	return ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *Quote) Update(ctx *context.Context, m *models.Quote) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Updates(m).Error
}

func (t *Quote) Get(ctx *context.Context, id string) (*models.Quote, error) {
	var result models.Quote
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get quote.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *Quote) Delete(ctx *context.Context, id string) error {
	var result models.Quote
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete quote.", zap.Error(err))
		return err
	}

	return err
}

func (t *Quote) GetAll(ctx *context.Context, ids []string) ([]*models.Quote, error) {
	var result []*models.Quote
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get quotes.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get quotes.", zap.Error(err))
		return nil, err
	}

	return result, err
}
