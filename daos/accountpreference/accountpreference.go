package accountpreference

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IAccountPreference interface {
	Upsert(ctx *context.Context, m ...*models.AccountPreference) error
	Get(ctx *context.Context, accountId, prefType string) (*models.AccountPreference, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.AccountPreference, error)
	Delete(ctx *context.Context, id string) error
}

type AccountPreference struct {
}

func NewAccountPreference() IAccountPreference {
	return &AccountPreference{}
}

func (t *AccountPreference) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "account_preferences"
}

func (t *AccountPreference) Upsert(ctx *context.Context, m ...*models.AccountPreference) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
	if err != nil {
		ctx.Log.Error("failed to upsert account preferences", zap.Error(err))
	}
	return err
}

func (t *AccountPreference) Get(ctx *context.Context, accountId, prefType string) (*models.AccountPreference, error) {
	var result models.AccountPreference
	q := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("account_id = ?", accountId)
	if prefType != "" {
		q = q.Where("type = ?", prefType)
	}
	err := q.First(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get accountpreference.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *AccountPreference) Delete(ctx *context.Context, id string) error {
	var result models.AccountPreference
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete accountpreference.", zap.Error(err))
		return err
	}

	return err
}

func (t *AccountPreference) GetAll(ctx *context.Context, ids []string) ([]*models.AccountPreference, error) {
	var result []*models.AccountPreference
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get accountpreferences.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get accountpreferences.", zap.Error(err))
		return nil, err
	}

	return result, err
}
