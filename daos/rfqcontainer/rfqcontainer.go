package rfqcontainer

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IRfqContainer interface {
	Upsert(ctx *context.Context, m ...*models.RfqContainer) error
	Get(ctx *context.Context, id string) (*models.RfqContainer, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.RfqContainer, error)
	Delete(ctx *context.Context, id string) error
	GetContainersByRfqId(ctx *context.Context, rfqID string) ([]*models.RfqContainer, error)
	Update(ctx *context.Context, m *models.RfqContainer) error
	GetByRfqIds(ctx *context.Context, rfqIds []string) ([]*models.RfqContainer, error)
	GetCountByRfq(ctx *context.Context, rfqId string) int64
}

type RfqContainer struct {
}

func NewRfqContainer() IRfqContainer {
	return &RfqContainer{}
}

func (t *RfqContainer) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "rfqs_containers"
}

func (t *RfqContainer) Upsert(ctx *context.Context, m ...*models.RfqContainer) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}
func (t *RfqContainer) Update(ctx *context.Context, m *models.RfqContainer) error {
	return ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("rfq_id = ?", m.RfqID).Updates(m).Error
}

func (t *RfqContainer) Get(ctx *context.Context, id string) (*models.RfqContainer, error) {
	var result models.RfqContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqcontainer.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *RfqContainer) Delete(ctx *context.Context, id string) error {
	var result models.RfqContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete rfqcontainer.", zap.Error(err))
	}

	return err
}

func (t *RfqContainer) GetAll(ctx *context.Context, ids []string) ([]*models.RfqContainer, error) {
	var result []*models.RfqContainer
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get rfqcontainers.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqcontainers.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *RfqContainer) GetContainersByRfqId(ctx *context.Context, rfqID string) ([]*models.RfqContainer, error) {
	var result []*models.RfqContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).
		Table(t.getTable(ctx)).
		Where("rfq_id = ? AND is_active = ?", rfqID, true).
		Find(&result).Error
	if err != nil {
		ctx.Log.Error("failed to get rfqcontainer by rfq id", zap.Error(err), zap.Any("rfq_id", rfqID))
		return nil, err
	}

	return result, nil
}

func (t *RfqContainer) GetByRfqIds(ctx *context.Context, rfqIds []string) ([]*models.RfqContainer, error) {
	var result []*models.RfqContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result, "rfq_id IN (?) AND is_deleted = false AND is_active = true", rfqIds).Error
	if err != nil {
		ctx.Log.Error("failed to get rfq containers by rfq ids", zap.Error(err), zap.Any("rfq_ids", rfqIds))
		return nil, err
	}

	return result, err
}

func (t *RfqContainer) GetCountByRfq(ctx *context.Context, rfqId string) int64 {

	var count int64
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("rfq_id = ?", rfqId).Count(&count).Error
	if err != nil {
		ctx.Log.Error("unable to get rfq containers count", zap.Error(err))
		return 0
	}

	return count
}
