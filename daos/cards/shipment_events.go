package workflow

import (
	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IShipmentEvents interface {
	Upsert(ctx *context.Context, m *models.ShipmentEvents) error
	GetAll(ctx *context.Context, req dtos.ShipmentEventsReq) ([]*models.ShipmentEvents, error)
}

type ShipmentEvents struct {
}

func NewWorkflows() IShipmentEvents {
	return &ShipmentEvents{}
}

func (t *ShipmentEvents) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".shipment_events"
}

func (t *ShipmentEvents) Upsert(ctx *context.Context, m *models.ShipmentEvents) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *ShipmentEvents) GetAll(ctx *context.Context, req dtos.ShipmentEventsReq) ([]*models.ShipmentEvents, error) {
	var result []*models.ShipmentEvents
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}
