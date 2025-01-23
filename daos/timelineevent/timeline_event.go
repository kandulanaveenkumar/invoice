package timelineevent

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type ITimelineEvent interface {
	Get(ctx *context.Context, id string) ([]*models.TimelineEvent, error)
	Upsert(ctx *context.Context, m ...*models.TimelineEvent) error
	GetByShipmentID(c *context.Context, shipmentID string) ([]*models.TimelineEvent, error)
	GetTimeLineEvent(ctx *context.Context, req *models.Card) ([]*models.TimelineEvent, error)
}

type TimelineEvent struct {
}

func NewTimelineEvent() ITimelineEvent {
	return &TimelineEvent{}
}

func (t *TimelineEvent) getTable(ctx *context.Context) string {
	if ctx.TenantID == "" {
		ctx.TenantID = "public"
	}

	return ctx.TenantID + "." + "timeline_event"
}

func (t *TimelineEvent) Get(ctx *context.Context, id string) ([]*models.TimelineEvent, error) {
	var result []*models.TimelineEvent
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id = ?", id).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get stock.", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (t *TimelineEvent) Upsert(ctx *context.Context, m ...*models.TimelineEvent) error {
	return ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *TimelineEvent) GetByShipmentID(ctx *context.Context,
	shipmentID string) ([]*models.TimelineEvent, error) {
	var result []*models.TimelineEvent
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Find(&result).Where("shipment_id = ?", shipmentID).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *TimelineEvent) GetTimeLineEvent(ctx *context.Context, req *models.Card) ([]*models.TimelineEvent, error) {
	var result []*models.TimelineEvent

	q := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if req.InstanceType == "shipment" {
		q = q.Where("shipment_id = ? ", req.InstanceId)
	}

	if req.InstanceType == "rfq" {
		q = q.Where("rfq_id = ? ", req.InstanceId)
	}

	if req.InstanceType == "quote" {
		q = q.Where("quote_id = ? ", req.InstanceId)
	}

	err := q.Find(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
