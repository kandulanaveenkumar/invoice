package cardaudits

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type ICardAudits interface {
	Upsert(ctx *context.Context, m ...*models.CardAudits) error
	Get(ctx *context.Context, id string) (*models.CardAudits, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.CardAudits, error)
	Delete(ctx *context.Context, id string) error
}

type CardAudits struct {
}

func NewCardAudits() ICardAudits {
	return &CardAudits{}
}

func (t *CardAudits) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "card_audits"
}

func (t *CardAudits) Upsert(ctx *context.Context, m ...*models.CardAudits) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Clauses(clause.OnConflict{UpdateAll: true}).Error
}

func (t *CardAudits) Get(ctx *context.Context, id string) (*models.CardAudits, error) {
	var result models.CardAudits
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get card_audits.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *CardAudits) Delete(ctx *context.Context, id string) error {
	var result models.CardAudits
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete card_audits.", zap.Error(err))
		return err
	}

	return err
}

func (t *CardAudits) GetAll(ctx *context.Context, ids []string) ([]*models.CardAudits, error) {
	var result []*models.CardAudits
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get card_auditss.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get card_auditss.", zap.Error(err))
		return nil, err
	}

	return result, err
}
