package airwaybillroute

import (
	"database/sql"
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IAirwayBillRoute interface {
	Upsert(ctx *context.Context, awbRoute *models.AirwayBillRoute, by uuid.UUID) (*models.AirwayBillRoute, error)
	UpsertAll(ctx *context.Context, awbRoutes []*models.AirwayBillRoute, by uuid.UUID) ([]*models.AirwayBillRoute, error)
	GetAll(ctx *context.Context, awbInfoId uuid.UUID, query string) ([]*models.AirwayBillRoute, error)
	Get(ctx *context.Context, id uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillRoute, error)
	Delete(ctx *context.Context, awbRoute *models.AirwayBillRoute, by uuid.UUID) error
}

type AirwayBillRoute struct {
}

func NewAirwayBillRoute() IAirwayBillRoute {
	return &AirwayBillRoute{}
}

func (t *AirwayBillRoute) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "awb_routes"
}

func (t *AirwayBillRoute) Upsert(ctx *context.Context, awbRoute *models.AirwayBillRoute, by uuid.UUID) (*models.AirwayBillRoute, error) {

	if awbRoute.Id == uuid.Nil {
		awbRoute.Id = uuid.New()
		awbRoute.CreatedAt = time.Now().UTC()
		awbRoute.CreatedBy = by.String()

	}
	awbRoute.UpdatedAt = time.Now().UTC()
	awbRoute.UpdatedBy = by.String()

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbRoute).Error
	if err != nil {
		ctx.Log.Error("unable to upsert data ", zap.Error(err))
		return nil, err
	}

	return awbRoute, nil
}

func (t *AirwayBillRoute) Get(ctx *context.Context, id uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillRoute, error) {
	awbRoute := &models.AirwayBillRoute{}

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
	err := tx.First(&awbRoute).Error
	if err != nil {
		return nil, err
	}

	return awbRoute, nil
}

func (t *AirwayBillRoute) Delete(ctx *context.Context, awbRoute *models.AirwayBillRoute, by uuid.UUID) error {

	awbRoute.DeletedAt = sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	awbRoute.DeletedBy = by.String()
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbRoute).Error
	if err != nil {
		ctx.Log.Error("unable to upsert data ", zap.Error(err))
		return err
	}
	err = ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(awbRoute).Error
	return err
}

func (t *AirwayBillRoute) GetAll(ctx *context.Context, awbInfoId uuid.UUID, query string) ([]*models.AirwayBillRoute, error) {
	awbRoutes := []*models.AirwayBillRoute{}

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if awbInfoId != uuid.Nil {
		tx.Where("awb_info_id = ?", awbInfoId)
	}
	if query != "" {
		tx.Where(query)
	}

	tx.Order("stop_number")
	err := tx.Find(&awbRoutes).Error
	if err != nil {
		return nil, err
	}

	return awbRoutes, nil
}

func (t *AirwayBillRoute) UpsertAll(ctx *context.Context, awbRoutes []*models.AirwayBillRoute, by uuid.UUID) ([]*models.AirwayBillRoute, error) {

	for _, awbRoute := range awbRoutes {

		if awbRoute.Id == uuid.Nil {

			awbRoute.Id = uuid.New()
			awbRoute.CreatedAt = time.Now().UTC()
			awbRoute.CreatedBy = by.String()

		}

		awbRoute.UpdatedAt = time.Now().UTC()
		awbRoute.UpdatedBy = by.String()
	}

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbRoutes).Error
	if err != nil {
		ctx.Log.Error("unable to save routes", zap.Error(err))
		return nil, err
	}

	return awbRoutes, nil
}
