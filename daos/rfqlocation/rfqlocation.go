package rfqlocation

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IRfqLocation interface {
	Upsert(ctx *context.Context, m ...*models.RfqLocation) error
	Get(ctx *context.Context, id string) (*models.RfqLocation, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.RfqLocation, error)
	GetByRfqId(ctx *context.Context, rfqID, locType string) ([]*models.RfqLocation, error)
	Delete(ctx *context.Context, id string) error
	DeleteByRfqID(ctx *context.Context, rfqID string) error
	Update(ctx *context.Context, m *models.RfqLocation) error
}

type RfqLocation struct {
}

func NewRfqLocation() IRfqLocation {
	return &RfqLocation{}
}

func (t *RfqLocation) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "rfq_locations"
}

func (t *RfqLocation) Upsert(ctx *context.Context, m ...*models.RfqLocation) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}
func (t *RfqLocation) Update(ctx *context.Context, m *models.RfqLocation) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("rfq_id = ?", m.RfqID).Updates(m).Error
	if err != nil {
		ctx.Log.Error("failed to update rfqlocation", zap.Error(err), zap.Any("rfq_id", m.RfqID))
	}
	return err
}

func (t *RfqLocation) Get(ctx *context.Context, id string) (*models.RfqLocation, error) {
	var result models.RfqLocation
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqlocation.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *RfqLocation) Delete(ctx *context.Context, id string) error {
	var result models.RfqLocation
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete rfqlocation.", zap.Error(err))
	}

	return err
}
func (t *RfqLocation) DeleteByRfqID(ctx *context.Context, rfqID string) error {
	var result models.RfqLocation
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "rfq_id = ?", rfqID).Error
	if err != nil {
		ctx.Log.Error("failed to delete rfq location by rfq id", zap.Error(err), zap.Any("rfq_id", rfqID))
	}

	return err
}

func (t *RfqLocation) GetAll(ctx *context.Context, ids []string) ([]*models.RfqLocation, error) {
	var result []*models.RfqLocation
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get rfqlocations.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqlocations.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *RfqLocation) GetByRfqId(ctx *context.Context, rfqID, locType string) ([]*models.RfqLocation, error) {
	var result []*models.RfqLocation
	q := ctx.DB.WithContext(ctx.Request.Context()).Debug().Table(t.getTable(ctx)).Where("rfq_id = ?", rfqID)
	if locType != "" {
		q = q.Where("type = ?", locType)
	}
	err := q.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqlocations.", zap.Error(err))
		return nil, err
	}

	return result, err

}
