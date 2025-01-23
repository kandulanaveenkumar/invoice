package rfqproduct

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IRfqProduct interface {
	Upsert(ctx *context.Context, m ...*models.RfqProduct) error
	Get(ctx *context.Context, id string) (*models.RfqProduct, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.RfqProduct, error)
	Delete(ctx *context.Context, id string) error
	GetRfqProductsByRfqId(ctx *context.Context, rfqID string) ([]*models.RfqProduct, error)
	Update(ctx *context.Context, m *models.RfqProduct) error
	GetByRfqIds(ctx *context.Context, rfqIds []string) ([]*models.RfqProduct, error)
}

type RfqProduct struct {
}

func NewRfqProduct() IRfqProduct {
	return &RfqProduct{}
}

func (t *RfqProduct) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "rfq_products"
}

func (t *RfqProduct) Upsert(ctx *context.Context, m ...*models.RfqProduct) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *RfqProduct) Get(ctx *context.Context, id string) (*models.RfqProduct, error) {
	var result models.RfqProduct
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ? AND is_active = true", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqproduct.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *RfqProduct) Delete(ctx *context.Context, id string) error {
	var result models.RfqProduct
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete rfqproduct.", zap.Error(err))
		return err
	}

	return err
}

func (t *RfqProduct) GetAll(ctx *context.Context, ids []string) ([]*models.RfqProduct, error) {
	var result []*models.RfqProduct
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get rfqproducts.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqproducts.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *RfqProduct) GetRfqProductsByRfqId(ctx *context.Context, rfqID string) ([]*models.RfqProduct, error) {
	var result []*models.RfqProduct
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("rfq_id = ? and is_active = ?", rfqID, true).Find(&result).Error
	if err != nil {
		ctx.Log.Error("failed to get rfqcontainer by rfq id", zap.Error(err), zap.Any("rfq_id", rfqID))
		return nil, err
	}
	return result, nil
}

func (t *RfqProduct) Update(ctx *context.Context, m *models.RfqProduct) error {
	return ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("rfq_id = ?", m.RfqID).Updates(m).Error
}

func (t *RfqProduct) GetByRfqIds(ctx *context.Context, rfqIds []string) ([]*models.RfqProduct, error) {
	var result []*models.RfqProduct
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result, "rfq_id IN (?) AND is_active = true", rfqIds).Error
	if err != nil {
		ctx.Log.Error("failed to get rfq products by rfq ids", zap.Error(err), zap.Any("rfq_ids", rfqIds))
		return nil, err
	}

	return result, err
}
