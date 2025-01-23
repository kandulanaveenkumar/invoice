package consolcontainer

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Iconsolcontainer interface {
	Get(ctx *context.Context, id string) (*models.ConsolContainer, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.ConsolContainer, error)
	Upsert(ctx *context.Context, m ...*models.ConsolContainer) error
	UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.ConsolContainer) error
	Delete(ctx *context.Context, id string) error
}

type consolcontainer struct {
}

func Newconsolcontainer() Iconsolcontainer {
	return &consolcontainer{}
}

func (t *consolcontainer) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "consol_containers"
}

func (t *consolcontainer) Upsert(ctx *context.Context, m ...*models.ConsolContainer) error {
	return ctx.DB.Table(t.getTable(ctx)).Save(m).Error
}

func (t *consolcontainer) UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.ConsolContainer) error {
	return tx.Table(t.getTable(ctx)).Save(m).Error
}

func (t *consolcontainer) Get(ctx *context.Context, id string) (*models.ConsolContainer, error) {
	var result models.ConsolContainer
	err := ctx.DB.Table(t.getTable(ctx)).First(&result, "consol_id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get consol container", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *consolcontainer) GetAll(ctx *context.Context, ids []string) ([]*models.ConsolContainer, error) {
	var result []*models.ConsolContainer
	if len(ids) == 0 {
		err := ctx.DB.Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get consol containers", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get consol containers", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *consolcontainer) Delete(ctx *context.Context, id string) error {
	var result models.ConsolContainer
	err := ctx.DB.Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete consol container", zap.Error(err))
		return err
	}

	return err
}
