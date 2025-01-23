package carddepartmentowner

import (
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type ICardDepartmentOwner interface {
	Upsert(ctx *context.Context, m ...*models.CardDepartmentOwner) error
	Get(ctx *context.Context, id string) (*models.CardDepartmentOwner, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.CardDepartmentOwner, error)
	Delete(ctx *context.Context, id string) error
	GetCardDepartmentOwnersWithFilter(ctx *context.Context, filter *models.CardDepartmentOwner) ([]*models.CardDepartmentOwner, error)
	UpsertCardDepartmentOwnersWithFilter(ctx *context.Context, instanceId, prevAssigned, assignedTo string) error
}

type CardDepartmentOwner struct {
}

func NewCardDepartmentOwner() ICardDepartmentOwner {
	return &CardDepartmentOwner{}
}

func (t *CardDepartmentOwner) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "card_department_owners"
}

func (t *CardDepartmentOwner) Upsert(ctx *context.Context, m ...*models.CardDepartmentOwner) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *CardDepartmentOwner) Get(ctx *context.Context, id string) (*models.CardDepartmentOwner, error) {
	var result models.CardDepartmentOwner
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get card_department_owner.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *CardDepartmentOwner) Delete(ctx *context.Context, id string) error {
	var result models.CardDepartmentOwner
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete card_department_owner.", zap.Error(err))
		return err
	}

	return err
}

func (t *CardDepartmentOwner) GetAll(ctx *context.Context, ids []string) ([]*models.CardDepartmentOwner, error) {
	var result []*models.CardDepartmentOwner
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get card_department_owners.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get card_department_owners.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *CardDepartmentOwner) GetCardDepartmentOwnersWithFilter(ctx *context.Context, filter *models.CardDepartmentOwner) ([]*models.CardDepartmentOwner, error) {
	var results []*models.CardDepartmentOwner

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	tx.Where(&filter).Order("created_at DESC")

	err := tx.Find(&results).Error
	if err != nil {
		ctx.Log.Error("error while fetching filtering card_department_owners", zap.Any("filter", filter), zap.Error(err))
		return nil, err
	}
	return results, nil
}

func (t *CardDepartmentOwner) UpsertCardDepartmentOwnersWithFilter(ctx *context.Context, instanceId, prevAssigned, assignedTo string) error {

	newValues := map[string]interface{}{
		"executive_id": assignedTo,
		"updated_at":   time.Now().UTC(),
	}

	return ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("instance_id = ? AND executive_id = ?", instanceId, prevAssigned).Updates(newValues).Error
}
