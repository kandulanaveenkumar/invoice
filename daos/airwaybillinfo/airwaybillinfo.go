package airwaybillinfo

import (
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IAirwayBillInfo interface {
	Upsert(ctx *context.Context, awbInfo *models.AirwayBillInfo, by uuid.UUID) (*models.AirwayBillInfo, error)
	GetAll(ctx *context.Context, shipmentId uuid.UUID, billType string, query string) ([]*models.AirwayBillInfo, error)
	Get(ctx *context.Context, id uuid.UUID, query string) (*models.AirwayBillInfo, error)
	GetForGenerateHAWBNumber(ctx *context.Context, portId string) (string, error)
	GetMAWBByShipmentId(ctx *context.Context, sids []string) ([]*models.DSRAWB, error)
	GetHAWBByShipmentId(ctx *context.Context, sids []string) ([]*models.DSRAWB, error)
}
type AirwayBillInfo struct {
}

func NewAirwayBillInfo() IAirwayBillInfo {
	return &AirwayBillInfo{}
}

func (t *AirwayBillInfo) getTable(ctx *context.Context) string {
	ctx.TenantID = "public"
	return ctx.TenantID + "." + "awb_info"
}

func (t *AirwayBillInfo) Upsert(ctx *context.Context, awbInfo *models.AirwayBillInfo, by uuid.UUID) (*models.AirwayBillInfo, error) {

	if awbInfo.Id == uuid.Nil {
		awbInfo.Id = uuid.New()
		awbInfo.CreatedAt = time.Now().UTC()
		awbInfo.CreatedBy = by

	}
	awbInfo.UpdatedAt = time.Now().UTC()
	awbInfo.UpdatedBy = by

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(&awbInfo).Error
	if err != nil {
		return nil, err
	}

	return awbInfo, nil
}

func (t *AirwayBillInfo) Get(ctx *context.Context, id uuid.UUID, query string) (*models.AirwayBillInfo, error) {
	awbInfo := &models.AirwayBillInfo{}

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if id != uuid.Nil {
		tx.Where("id = ?", id)
	}
	if query != "" {
		tx.Where(query)
	}
	err := tx.First(&awbInfo).Error
	if err != nil {
		return nil, err
	}

	return awbInfo, nil
}

func (t *AirwayBillInfo) GetAll(ctx *context.Context, shipmentId uuid.UUID, billType string, query string) ([]*models.AirwayBillInfo, error) {
	awbsInfo := []*models.AirwayBillInfo{}

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	if shipmentId != uuid.Nil {
		tx.Joins("awb_house ON awb_info.id = awb_house.awb_info_id AND awb_house.shipment_id = ?", shipmentId)
	}
	if billType != "" {
		tx.Where("awb_info.type = ?", billType)
	}
	if query != "" {
		tx.Where(query)
	}

	tx.Order("created_at")
	err := tx.Find(&awbsInfo).Error
	if err != nil {
		return nil, err
	}

	return awbsInfo, nil
}

func (t *AirwayBillInfo) GetForGenerateHAWBNumber(ctx *context.Context, portId string) (string, error) {

	awbInfo := &models.AirwayBillInfo{}
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	tx.Where("type = ? AND issuer_port_code = ?", globals.HouseAirwayBill, portId)
	tx.Order("created_at DESC, number DESC")
	err := tx.First(&awbInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	return awbInfo.Number, nil
}

func (t *AirwayBillInfo) GetMAWBByShipmentId(ctx *context.Context, sids []string) ([]*models.DSRAWB, error) {
	awb_info := t.getTable(ctx)
	awb_master := ctx.TenantID + "." + "awb_master"
	mawb := []*models.DSRAWB{}
	err := ctx.DB.WithContext(ctx.Request.Context()).Raw(`Select m.shipment_id,i.number from `+awb_info+` i INNER JOIN `+awb_master+` m on i.id=m.awb_info_id where m.shipment_id IN ?`, sids).Scan(&mawb).Error

	if err != nil {
		ctx.Log.Error("Error in fetching mawb", zap.Error(err))
		return nil, err
	}
	return mawb, nil
}

func (t *AirwayBillInfo) GetHAWBByShipmentId(ctx *context.Context, sids []string) ([]*models.DSRAWB, error) {
	awb_house := ctx.TenantID + "." + "awb_house"
	awb_info := t.getTable(ctx)
	hawb := []*models.DSRAWB{}
	err := ctx.DB.WithContext(ctx.Request.Context()).Raw(`Select h.shipment_id,i.number from `+awb_info+` i INNER JOIN `+awb_house+` h on i.id=h.awb_info_id where h.shipment_id IN ?`, sids).Scan(&hawb).Error
	if err != nil {
		ctx.Log.Error("Error in fetching hawb", zap.Error(err))
		return nil, err
	}
	return hawb, nil
}
