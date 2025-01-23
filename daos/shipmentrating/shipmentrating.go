package shipmentrating

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IShipmentRating interface {
	Upsert(ctx *context.Context, m ...*models.ShipmentRating) error
	Get(ctx *context.Context, id string) (*models.ShipmentRating, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.ShipmentRating, error)
	Delete(ctx *context.Context, id string) error
	GetByShipment(ctx *context.Context, shipmentIds []string) ([]*models.ShipmentRating, error)
	GetRatingByShipmentId(ctx *context.Context, sid string) (*models.ShipmentRating, error)
	GetShipmentRatingActivities(ctx *context.Context, cid string) ([]*models.Shipment, error)
}

type ShipmentRating struct {
}

func NewShipmentRating() IShipmentRating {
	return &ShipmentRating{}
}

func (t *ShipmentRating) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "shipment_ratings"
}

func (t *ShipmentRating) getShipmentTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "shipments"
}

func (t *ShipmentRating) Upsert(ctx *context.Context, m ...*models.ShipmentRating) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *ShipmentRating) Get(ctx *context.Context, id string) (*models.ShipmentRating, error) {
	var result models.ShipmentRating
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentrating.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *ShipmentRating) Delete(ctx *context.Context, id string) error {
	var result models.ShipmentRating
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete shipmentrating.", zap.Error(err))
		return err
	}

	return err
}

func (t *ShipmentRating) GetAll(ctx *context.Context, ids []string) ([]*models.ShipmentRating, error) {
	var result []*models.ShipmentRating
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get shipmentratings.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentratings.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *ShipmentRating) GetByShipment(ctx *context.Context, shipmentIds []string) ([]*models.ShipmentRating, error) {
	var result []*models.ShipmentRating
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result, "shipment_id IN (?)", shipmentIds).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentrating.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *ShipmentRating) GetRatingByShipmentId(ctx *context.Context, sid string) (*models.ShipmentRating, error) {
	result := &models.ShipmentRating{}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "shipment_id = ?", sid).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentrating.", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (t *ShipmentRating) GetShipmentRatingActivities(ctx *context.Context, cid string) ([]*models.Shipment, error) {
	var result []*models.Shipment
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getShipmentTable(ctx)+" s").
		Select("s.id, s.code, s.pol, s.pod, s.is_door_pickup,s.is_door_delivery,s.type").
		Joins("LEFT JOIN  shipment_ratings  as sr ON s.id = sr.shipment_id").
		Where("s.company_id = ?", cid).
		Where("s.status = ?", constants.ShipmentCompleted).
		Where("s.is_deleted = ? ", false).
		Where("sr.ratings=0 or sr.ratings is null")

	err := query.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Error in fetching rfq card details", zap.Error(err))
		return nil, err
	}

	return result, nil
}
