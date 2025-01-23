package airwaybilldoc

import (
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IAirwayBillDoc interface {
	Upsert(ctx *context.Context, awbDocs *models.AirwayBillDocs, by uuid.UUID) (*models.AirwayBillDocs, error)
	UpsertAll(ctx *context.Context, awbDocs []*models.AirwayBillDocs, by uuid.UUID) ([]*models.AirwayBillDocs, error)
	GetAll(ctx *context.Context, awbInfoId uuid.UUID, query string) ([]*models.AirwayBillDocs, error)
	Get(ctx *context.Context, id uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillDocs, error)
}

type AirwayBillDoc struct {
}

func NewAirwayBillDoc() IAirwayBillDoc {
	return &AirwayBillDoc{}
}

func (t *AirwayBillDoc) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "awb_docs"
}

func (t *AirwayBillDoc) Upsert(ctx *context.Context, awbDoc *models.AirwayBillDocs, by uuid.UUID) (*models.AirwayBillDocs, error) {

	currentTime := time.Now().UTC()
	if awbDoc.Id == uuid.Nil {
		awbDoc.Id = uuid.New()
		awbDoc.CreatedAt = currentTime
		awbDoc.CreatedBy = by

	}
	awbDoc.UpdatedAt = currentTime
	awbDoc.UpdatedBy = by

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbDoc).Error
	if err != nil {
		ctx.Log.Error("unable to upsert docs", zap.Error(err))
		return nil, err
	}

	return awbDoc, nil
}

func (t *AirwayBillDoc) Get(ctx *context.Context, id uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillDocs, error) {
	awbDoc := &models.AirwayBillDocs{}

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
	err := tx.First(&awbDoc).Error
	if err != nil {
		return nil, err
	}

	return awbDoc, nil
}

func (t *AirwayBillDoc) GetAll(ctx *context.Context, awbInfoId uuid.UUID, query string) ([]*models.AirwayBillDocs, error) {
	awbDocs := []*models.AirwayBillDocs{}

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if awbInfoId != uuid.Nil {
		tx.Where("awb_info_id = ?", awbInfoId)
	}
	if query != "" {
		tx.Where(query)
	}

	tx.Order("created_at")
	err := tx.Find(&awbDocs).Error
	if err != nil {
		return nil, err
	}

	return awbDocs, nil
}

func (t *AirwayBillDoc) UpsertAll(ctx *context.Context, awbDocs []*models.AirwayBillDocs, by uuid.UUID) ([]*models.AirwayBillDocs, error) {

	currentTime := time.Now().UTC()

	for _, awbDoc := range awbDocs {

		if awbDoc.Id == uuid.Nil {
			awbDoc.Id = uuid.New()
			awbDoc.CreatedAt = currentTime
			awbDoc.CreatedBy = by

		}

		awbDoc.UpdatedAt = currentTime
		awbDoc.UpdatedBy = by
	}

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbDocs).Error
	if err != nil {
		ctx.Log.Error("unable to upsert all", zap.Error(err))
		return nil, err
	}

	return awbDocs, nil
}
