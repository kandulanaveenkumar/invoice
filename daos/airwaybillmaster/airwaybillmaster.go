package airwaybillmaster

import (
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type IAirwayBillMaster interface {
	Upsert(ctx *context.Context, awbMaster *models.AirwayBillMaster, by uuid.UUID) (*models.AirwayBillMaster, error)
	UpsertAll(ctx *context.Context, awbMasters []*models.AirwayBillMaster, by uuid.UUID) ([]*models.AirwayBillMaster, error)
	GetAll(ctx *context.Context, awbInfoId uuid.UUID, pol string, pod string, linerId string, query string) ([]*models.AirwayBillMaster, error)
	Get(ctx *context.Context, id uuid.UUID, shipmentId uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillMaster, error)
	GetMastersWithStockNumbers(ctx *context.Context, stockNumbers []string, query string) ([]*models.MawbStockWithBookings, error)
}

type AirwayBillMaster struct {
}

func NewAirwayBillMaster() IAirwayBillMaster {
	return &AirwayBillMaster{}
}

func (t *AirwayBillMaster) getTable(ctx *context.Context) string {
	if ctx.TenantID == "" {
		ctx.TenantID = "public"
	}
	return ctx.TenantID + "." + "awb_master"
}

func (t *AirwayBillMaster) Upsert(ctx *context.Context, awbMaster *models.AirwayBillMaster, by uuid.UUID) (*models.AirwayBillMaster, error) {

	if awbMaster.Id == uuid.Nil {
		awbMaster.Id = uuid.New()
		awbMaster.CreatedAt = time.Now().UTC()
		awbMaster.CreatedBy = by

	}
	awbMaster.UpdatedAt = time.Now().UTC()
	awbMaster.UpdatedBy = by

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbMaster).Clauses(clause.OnConflict{DoNothing: true}).Error
	if err != nil {
		return nil, err
	}

	return awbMaster, nil

}

func (t *AirwayBillMaster) UpsertAll(ctx *context.Context, awbMasters []*models.AirwayBillMaster, by uuid.UUID) ([]*models.AirwayBillMaster, error) {

	currentTime := time.Now().UTC()

	for _, awbMaster := range awbMasters {
		if awbMaster.Id == uuid.Nil {
			awbMaster.Id = uuid.New()
			awbMaster.CreatedAt = currentTime
			awbMaster.CreatedBy = by
		}

		awbMaster.UpdatedAt = currentTime
		awbMaster.UpdatedBy = by
	}

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbMasters).Error
	if err != nil {
		ctx.Log.Error("Failed to upsert AWB masters", zap.Error(err))
		return nil, err
	}

	return awbMasters, nil
}

func (t *AirwayBillMaster) Get(ctx *context.Context, id uuid.UUID, shipmentId uuid.UUID, awbInfoId uuid.UUID, query string) (*models.AirwayBillMaster, error) {
	awbMaster := &models.AirwayBillMaster{}

	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
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
	err := tx.First(&awbMaster).Error
	if err != nil {
		return nil, err
	}

	return awbMaster, nil
}

func (t *AirwayBillMaster) GetAll(ctx *context.Context, awbInfoId uuid.UUID, pol string, pod string, linerId string, query string) ([]*models.AirwayBillMaster, error) {

	awbMasters := []*models.AirwayBillMaster{}
	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if awbInfoId != uuid.Nil {
		tx.Where("awb_info_id = ?", awbInfoId)
	}
	if pol != "" {
		tx.Where("pol = ?", pol)
	}
	if pod != "" {
		tx.Where("pod = ?", pod)
	}
	if linerId != "" {
		tx.Where("liner_code = ?", linerId)
	}
	if query != "" {
		tx.Where(query)
	}

	tx.Order("created_at")
	err := tx.Find(&awbMasters).Error
	if err != nil {
		return nil, err
	}

	return awbMasters, nil
}

func (t *AirwayBillMaster) GetMastersWithStockNumbers(ctx *context.Context, stockNumbers []string, query string) ([]*models.MawbStockWithBookings, error) {
	awbMasters := []*models.MawbStockWithBookings{}
	querystr := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("awb_master.awb_info_id , awb_info.number, array_agg(awb_master.shipment_id::text order by awb_master.created_at asc ) as shipment_ids").Joins("JOIN awb_info ON awb_info.id = awb_master.awb_info_id")
	if len(stockNumbers) > 0 {
		querystr.Where(" awb_info.number IN ?", stockNumbers)
	}
	if query != "" {
		querystr.Where(query)
	}
	querystr.Group("awb_info.number, awb_master.awb_info_id, awb_info.created_at")
	querystr.Order("awb_info.created_at")
	err := querystr.Find(&awbMasters).Error
	if err != nil {
		return nil, err
	}

	return awbMasters, nil
}
