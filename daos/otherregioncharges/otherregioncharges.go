package otherregioncharges

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IOtherRegionCharges interface {
	Upsert(ctx *context.Context, m ...*models.OtherRegionCharges) error
	Get(ctx *context.Context, id string) (*models.OtherRegionCharges, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.OtherRegionCharges, error)
	Delete(ctx *context.Context, id string) error
	GetForRegions(ctx *context.Context, region, bookingRegion, product string) ([]*models.OtherRegionCharges, error)
}

type OtherRegionCharges struct {
}

func NewOtherRegionCharges() IOtherRegionCharges {
	return &OtherRegionCharges{}
}

func (t *OtherRegionCharges) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "other_region_charges"
}

func (t *OtherRegionCharges) Upsert(ctx *context.Context, m ...*models.OtherRegionCharges) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *OtherRegionCharges) Get(ctx *context.Context, id string) (*models.OtherRegionCharges, error) {
	var result models.OtherRegionCharges
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get otherregioncharges.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *OtherRegionCharges) Delete(ctx *context.Context, id string) error {
	var result models.OtherRegionCharges
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete otherregioncharges.", zap.Error(err))
		return err
	}

	return err
}

func (t *OtherRegionCharges) GetAll(ctx *context.Context, ids []string) ([]*models.OtherRegionCharges, error) {
	var result []*models.OtherRegionCharges
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get otherregionchargess.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get otherregionchargess.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (h *OtherRegionCharges) GetForRegions(ctx *context.Context, region, bookingRegion, product string) ([]*models.OtherRegionCharges, error) {
	handlingCharges := []*models.OtherRegionCharges{}

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(h.getTable(ctx)).
		Where("region_id = ? OR region_id = 'All'", region).
		Where("booking_Region = ? OR booking_Region = 'All'", bookingRegion).
		Where("product = ? OR product = 'All'", product).
		Find(&handlingCharges).Error
	if err != nil {
		ctx.Log.Error("Unable to upsert handling charges.", zap.Error(err))
		return nil, err
	}

	return handlingCharges, nil
}
