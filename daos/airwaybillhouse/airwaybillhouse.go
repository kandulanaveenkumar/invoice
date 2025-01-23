package airwaybillhouse

import (
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type IAirwayBillHouse interface {
	Upsert(ctx *context.Context, awbHouse *models.AirwayBillHouse, by uuid.UUID) (*models.AirwayBillHouse, error)
	GetAll(ctx *context.Context, shipmentId uuid.UUID, query string) ([]*models.AirwayBillHouse, error)
	Get(ctx *context.Context, id uuid.UUID, shipmentId uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillHouse, error)
}

type AirwayBillHouse struct {
}

func NewAirwayBillHouse() IAirwayBillHouse {
	return &AirwayBillHouse{}
}

func (t *AirwayBillHouse) getTable(ctx *context.Context) string {
	ctx.TenantID = "public"
	return ctx.TenantID + "." + "awb_house"
}

func (t *AirwayBillHouse) Upsert(ctx *context.Context, awbHouse *models.AirwayBillHouse, by uuid.UUID) (*models.AirwayBillHouse, error) {

	currentTime := time.Now().UTC()
	if awbHouse.Id == uuid.Nil {
		awbHouse.Id = uuid.New()
		awbHouse.CreatedAt = currentTime
		awbHouse.CreatedBy = by

	}
	awbHouse.UpdatedAt = currentTime
	awbHouse.UpdatedBy = by

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbHouse).Clauses(clause.OnConflict{DoNothing: true}).Error
	if err != nil {
		return nil, err
	}

	return awbHouse, nil
}

func (t *AirwayBillHouse) Get(ctx *context.Context, id uuid.UUID, shipmentId uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillHouse, error) {
	awbHouse := &models.AirwayBillHouse{}

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if id != uuid.Nil {
		tx.Where("id = ?", id)
	}
	if shipmentId != uuid.Nil {
		tx.Where("shipment_id = ?", shipmentId)
	}
	if awbInfoId != uuid.Nil {
		tx.Where("awb_info_id = ?", awbInfoId)
	}
	if query != "" {
		tx.Where(query)
	}
	err := tx.First(&awbHouse).Error
	if err != nil {
		return nil, err
	}

	return awbHouse, nil
}

func (t *AirwayBillHouse) GetAll(ctx *context.Context, shipmentId uuid.UUID, query string) ([]*models.AirwayBillHouse, error) {
	awbHouses := []*models.AirwayBillHouse{}

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if shipmentId != uuid.Nil {
		tx.Where("shipment_id = ?", shipmentId)
	}
	if query != "" {
		tx.Where(query)
	}

	tx.Order("created_at")
	err := tx.Find(&awbHouses).Error
	if err != nil {
		return nil, err
	}

	return awbHouses, nil
}
