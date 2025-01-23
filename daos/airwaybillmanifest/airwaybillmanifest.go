package airwaybillmanifest

import (
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
)

type IAirwayBillManifest interface {
	Upsert(ctx *context.Context, awbManifest *models.AirwayBillManifest, by uuid.UUID) (*models.AirwayBillManifest, error)
	GetAll(ctx *context.Context, id uuid.UUID, mawbId uuid.UUID, hawbId uuid.UUID, query string) ([]*models.AirwayBillManifest, error)
	Get(ctx *context.Context, id uuid.UUID, mawbId uuid.UUID, hawbId uuid.UUID, query string) (*models.AirwayBillManifest, error)
}

type AirwayBillManifest struct {
}

func NewAirwayBillManifest() IAirwayBillManifest {
	return &AirwayBillManifest{}
}

func (t *AirwayBillManifest) getTable(ctx *context.Context) string {
	if ctx.TenantID == "" {
		ctx.TenantID = "public"
	}
	return ctx.TenantID + "." + "awb_manifest"
}

func (t *AirwayBillManifest) Upsert(ctx *context.Context, awbManifest *models.AirwayBillManifest, by uuid.UUID) (*models.AirwayBillManifest, error) {

	currentTime := time.Now().UTC()
	if awbManifest.Id == uuid.Nil {
		awbManifest.Id = uuid.New()
		awbManifest.CreatedAt = currentTime
		awbManifest.CreatedBy = by

	}
	awbManifest.UpdatedAt = currentTime
	awbManifest.UpdatedBy = by

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbManifest).Error
	if err != nil {
		return nil, err
	}

	return awbManifest, nil
}

func (t *AirwayBillManifest) Get(ctx *context.Context, id uuid.UUID, mawbId uuid.UUID, hawbId uuid.UUID, query string) (*models.AirwayBillManifest, error) {
	awbManifest := &models.AirwayBillManifest{}

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if id != uuid.Nil {
		tx.Where("id = ?", id)
	}
	if mawbId != uuid.Nil {
		tx.Where("mawb_id = ?", mawbId)
	}
	if hawbId != uuid.Nil {
		tx.Where("hawb_id = ?", hawbId)
	}
	if query != "" {
		tx.Where(query)
	}
	err := tx.First(&awbManifest).Error
	if err != nil {
		return nil, err
	}

	return awbManifest, nil
}

// func (t *AirwayBillManifest) Delete(ctx *context.Context, id string) error {
// 	var result models.AirwayBillManifest
// 	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
// 	if err != nil {
// 		ctx.Log.Error("Unable to delete airwaybillmanifest.", zap.Error(err))
// 		return err
// 	}

// 	return err
// }

func (t *AirwayBillManifest) GetAll(ctx *context.Context, id uuid.UUID, mawbId uuid.UUID, hawbId uuid.UUID, query string) ([]*models.AirwayBillManifest, error) {
	awbManifests := []*models.AirwayBillManifest{}

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if id != uuid.Nil {
		tx.Where("id = ?", id)
	}
	if mawbId != uuid.Nil {
		tx.Where("mawb_id = ?", mawbId)
	}
	if hawbId != uuid.Nil {
		tx.Where("hawb_id = ?", hawbId)
	}

	if query != "" {
		tx.Where(query)
	}

	tx.Order("created_at")
	err := tx.Find(&awbManifests).Error
	if err != nil {
		return nil, err
	}

	return awbManifests, nil
}
