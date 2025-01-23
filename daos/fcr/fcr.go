package fcr

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IFcr interface {
	Upsert(ctx *context.Context, m ...*models.Fcr) error
	Get(ctx *context.Context, blNo string, query string) (*models.Fcr, error)
	GetAll(ctx *context.Context, shipmentId string, query string) ([]*models.Fcr, error)
}

type Fcr struct {
}

func NewFcr() IFcr {
	return &Fcr{}
}

func (f *Fcr) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "fcr"
}

func (f *Fcr) Upsert(ctx *context.Context, m ...*models.Fcr) error {
	return ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(f.getTable(ctx)).Save(m).Error
}

func (f *Fcr) Get(ctx *context.Context, blNo string, query string) (*models.Fcr, error) {

	result := &models.Fcr{}

	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(f.getTable(ctx))

	if blNo != "" {
		tx.Where("bl_no = ?", blNo)
	}

	if query != "" {
		tx.Where(query)
	}

	tx.Order("created_at DESC")

	err := tx.First(&result).Error
	if err != nil {
		ctx.Log.Error("unable to get fcr details", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (f *Fcr) GetAll(ctx *context.Context, shipmentId string, query string) ([]*models.Fcr, error) {

	result := []*models.Fcr{}

	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(f.getTable(ctx))

	if shipmentId != "" {
		tx.Where("shipment_id = ?", shipmentId)
	}
	if query != "" {
		tx.Where(query)
	}
	tx.Order("created_at")

	err := tx.Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
