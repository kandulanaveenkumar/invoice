package airwaybilllabels

import (
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
)

type IAirwayBillLabels interface {
	Upsert(ctx *context.Context, awbHouse *models.AirwayBillLabels, by uuid.UUID) (*models.AirwayBillLabels, error)
	GetAll(ctx *context.Context, shipmentId uuid.UUID, query string) ([]*models.AirwayBillLabels, error)
	Get(ctx *context.Context, id uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillLabels, error)
}

type AirwayBillLabels struct {
}

func NewAirwayBillLabels() IAirwayBillLabels {
	return &AirwayBillLabels{}
}

func (t *AirwayBillLabels) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "awb_labels"
}

func (t *AirwayBillLabels) Upsert(ctx *context.Context, awbHouse *models.AirwayBillLabels, by uuid.UUID) (*models.AirwayBillLabels, error) {

	currentTime := time.Now().UTC()
	if awbHouse.Id == uuid.Nil {
		awbHouse.Id = uuid.New()
		awbHouse.CreatedAt = currentTime
		awbHouse.CreatedBy = by

	}
	awbHouse.UpdatedAt = currentTime
	awbHouse.UpdatedBy = by

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbHouse).Error
	if err != nil {
		return nil, err
	}

	return awbHouse, nil
}

func (t *AirwayBillLabels) Get(ctx *context.Context, id uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillLabels, error) {
	var awbHouse *models.AirwayBillLabels

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if id != uuid.Nil {
		tx.Where("id = ?", id)
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

// func (t *AirwayBillLabels) Delete(ctx *context.Context, id string) error {
// 	var result models.AirwayBillLabels
// 	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
// 	if err != nil {
// 		ctx.Log.Error("Unable to delete airwaybilllabels.", zap.Error(err))
// 		return err
// 	}

// 	return err
// }

func (t *AirwayBillLabels) GetAll(ctx *context.Context, id uuid.UUID, query string) ([]*models.AirwayBillLabels, error) {
	var awbHouses []*models.AirwayBillLabels

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
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
