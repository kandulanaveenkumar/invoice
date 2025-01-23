package airwaybillcharges

import (
	"database/sql"
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IAirwayBillCharges interface {
	Upsert(ctx *context.Context, awbCharge *models.AirwayBillCharge, by uuid.UUID) (*models.AirwayBillCharge, error)
	UpsertAll(ctx *context.Context, awbCharges []*models.AirwayBillCharge, by uuid.UUID) ([]*models.AirwayBillCharge, error)
	GetAll(ctx *context.Context, awbInfoId uuid.UUID, query string) ([]*models.AirwayBillCharge, error)
	Get(ctx *context.Context, id uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillCharge, error)
	Delete(ctx *context.Context, awbCharge *models.AirwayBillCharge, id uuid.UUID) error
}

type AirwayBillCharges struct {
}

func NewAirwayBillCharges() IAirwayBillCharges {
	return &AirwayBillCharges{}
}

func (t *AirwayBillCharges) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "awb_charges"
}

func (t *AirwayBillCharges) Upsert(ctx *context.Context, awbCharge *models.AirwayBillCharge, by uuid.UUID) (*models.AirwayBillCharge, error) {

	currentTime := time.Now().UTC()
	if awbCharge.Id == uuid.Nil {
		awbCharge.Id = uuid.New()
		awbCharge.CreatedAt = currentTime
		awbCharge.CreatedBy = by.String()

	}
	awbCharge.UpdatedAt = currentTime
	awbCharge.UpdatedBy = by.String()

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbCharge).Error
	if err != nil {
		ctx.Log.Error("unable to upsert charges", zap.Error(err))
		return nil, err
	}

	return awbCharge, nil

}

func (t *AirwayBillCharges) Get(ctx *context.Context, id uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillCharge, error) {
	awbCharge := &models.AirwayBillCharge{}

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
	err := tx.First(&awbCharge).Error
	if err != nil {
		return nil, err
	}

	return awbCharge, nil
}

func (t *AirwayBillCharges) Delete(ctx *context.Context, awbCharge *models.AirwayBillCharge, by uuid.UUID) error {

	awbCharge.DeletedAt = sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	awbCharge.DeletedBy = by.String()

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(awbCharge).Error
	if err != nil {
		ctx.Log.Error("unable to save updateby", zap.Error(err))
		return err
	}
	err = ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(awbCharge).Error
	return err
}

func (t *AirwayBillCharges) GetAll(ctx *context.Context, awbInfoId uuid.UUID, query string) ([]*models.AirwayBillCharge, error) {
	awbCharges := []*models.AirwayBillCharge{}

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(globals.TableAWBCharges)
	if awbInfoId != uuid.Nil {
		tx.Where("awb_info_id = ?", awbInfoId)
	}
	if query != "" {
		tx.Where(query)
	}

	tx.Order("created_at")
	err := tx.Find(&awbCharges).Error
	if err != nil {
		return nil, err
	}

	return awbCharges, nil
}

func (t *AirwayBillCharges) UpsertAll(ctx *context.Context, awbCharges []*models.AirwayBillCharge, by uuid.UUID) ([]*models.AirwayBillCharge, error) {

	currentTime := time.Now().UTC()

	for _, awbCharge := range awbCharges {

		if awbCharge.Id == uuid.Nil {
			awbCharge.Id = uuid.New()
			awbCharge.CreatedAt = currentTime
			awbCharge.CreatedBy = by.String()
		}

		awbCharge.UpdatedAt = currentTime
		awbCharge.UpdatedBy = by.String()
	}

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbCharges).Error
	if err != nil {
		ctx.Log.Error("error saving charges", zap.Error(err))
		return nil, err
	}

	return awbCharges, nil
}
