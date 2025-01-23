package arrivaldeliverydocs

import (
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IArrivalDeliveryDocs interface {
	Upsert(ctx *context.Context, addocs *models.ArrivalDeliveryDocs, by uuid.UUID) (*models.ArrivalDeliveryDocs, error)
	GetArrivalDeliveryDocsById(ctx *context.Context, id uuid.UUID, shipmentId uuid.UUID, HblNumber string) (*models.ArrivalDeliveryDocs, error)
	GetContainerInfoById(ctx *context.Context, id uuid.UUID) ([]*models.ContainerInfo, error)
	GetAirTranshipmentInfo(ctx *context.Context, id uuid.UUID) ([]*models.AirTranshipmentInfo, error)
	GetDoBlInfo(ctx *context.Context, id uuid.UUID, blno string) (*models.DoBlInfo, error)
	GetTranshipmentInfoById(ctx *context.Context, id uuid.UUID) ([]*models.TranshipmentInfo, error)
	DeleteAddocsTranshipmentInfo(ctx *context.Context, addocId uuid.UUID) error
	DeleteAddocsAirTranshipmentInfo(ctx *context.Context, addocId uuid.UUID) error
	DeleteAddocsContainerInfo(ctx *context.Context, addocId uuid.UUID) error
	GetMaxDoBlNumber(ctx *context.Context) int64
}

type ArrivalDeliveryDocs struct {
}

func NewArrivalDeliveryDocs() IArrivalDeliveryDocs {
	return &ArrivalDeliveryDocs{}
}

func (b *ArrivalDeliveryDocs) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + globals.TableADDocs
}

func (b *ArrivalDeliveryDocs) getAdDocContainerTable(ctx *context.Context) string {
	return ctx.TenantID + "." + globals.TableADDocsContainerInfo
}

func (b *ArrivalDeliveryDocs) getAdDocTranshipmentTable(ctx *context.Context) string {
	return ctx.TenantID + "." + globals.TableADDocsTranshipmentInfo
}

func (b *ArrivalDeliveryDocs) getAdDocAirTranshipmentTable(ctx *context.Context) string {
	return ctx.TenantID + "." + globals.TableADDocsAirTranshipmentInfo
}

func (b *ArrivalDeliveryDocs) getAdDocDoBlTable(ctx *context.Context) string {
	return ctx.TenantID + "." + globals.TableADDocsDoBlInfo
}

func (b *ArrivalDeliveryDocs) Upsert(ctx *context.Context, addocs *models.ArrivalDeliveryDocs, by uuid.UUID) (*models.ArrivalDeliveryDocs, error) {

	currentTime := time.Now().UTC().Unix()

	if addocs.Id == uuid.Nil {
		addocs.Id = uuid.New()
		addocs.CreatedAt = currentTime
		addocs.CreatedBy = by
	}

	addocs.UpdatedAt = currentTime
	addocs.UpdatedBy = by

	err := ctx.DB.Debug().Debug().WithContext(ctx.Request.Context()).Table(b.getTable(ctx)).Save(&addocs).Error
	if err != nil {
		return nil, err
	}

	containerinfo := []*models.ContainerInfo{}

	for _, value := range addocs.ContainerInfo {

		containerInfo := &models.ContainerInfo{}

		if value.Id == uuid.Nil {
			containerInfo.Id = uuid.New()
			containerInfo.ContainerNumber = value.ContainerNumber
			containerInfo.Seal = value.Seal
			containerInfo.Type = value.Type
			containerInfo.Weight = value.Weight
			containerInfo.Volume = value.Volume
			containerInfo.Packs = value.Packs
			containerInfo.Marks = value.Marks
			containerInfo.Metadata.CreatedAt = currentTime
			containerInfo.AddocID = addocs.Id
			containerInfo.Metadata.CreatedBy = by
		}

		containerInfo.Metadata.UpdatedAt = currentTime
		containerInfo.Metadata.UpdatedBy = by
		containerinfo = append(containerinfo, containerInfo)
	}

	if len(containerinfo) > 0 {
		err = ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocContainerTable(ctx)).Save(&containerinfo).Error
		if err != nil {
			return nil, err
		}
	}

	transhipmentinfo := []*models.TranshipmentInfo{}

	for _, value := range addocs.TranshipmentInfo {

		transhipmentInfo := &models.TranshipmentInfo{}

		if value.Id == uuid.Nil {
			transhipmentInfo.Id = uuid.New()
		} else {
			transhipmentInfo.Id = value.Id
		}

		transhipmentInfo.AddocID = addocs.Id
		transhipmentInfo.Pol = value.Pol
		transhipmentInfo.Pod = value.Pod
		transhipmentInfo.Vesselname = value.Vesselname
		transhipmentInfo.Voyageno = value.Voyageno
		transhipmentInfo.Etd = value.Etd
		transhipmentInfo.TransitDays = value.TransitDays
		transhipmentInfo.Metadata.CreatedAt = currentTime
		transhipmentInfo.Metadata.CreatedBy = by
		transhipmentInfo.Metadata.UpdatedAt = currentTime
		transhipmentInfo.Metadata.UpdatedBy = by
		transhipmentinfo = append(transhipmentinfo, transhipmentInfo)
	}

	if len(transhipmentinfo) > 0 {
		err = ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocTranshipmentTable(ctx)).Save(&transhipmentinfo).Error
		if err != nil {
			return nil, err
		}
	}

	airTranshipmentInfos := []*models.AirTranshipmentInfo{}

	for _, value := range addocs.AirTranshipmentInfo {

		airTranshipmentInfo := &models.AirTranshipmentInfo{}

		if value.Id == uuid.Nil {
			airTranshipmentInfo.Id = uuid.New()
		} else {
			airTranshipmentInfo.Id = value.Id
		}

		airTranshipmentInfo.AddocID = addocs.Id
		airTranshipmentInfo.Stops = value.Stops
		airTranshipmentInfo.Pod = value.Pod
		airTranshipmentInfo.Flight = value.Flight
		airTranshipmentInfo.Date = value.Date
		airTranshipmentInfo.From = value.From
		airTranshipmentInfo.To = value.To
		airTranshipmentInfo.Metadata.CreatedAt = currentTime
		airTranshipmentInfo.Metadata.CreatedBy = by
		airTranshipmentInfo.Metadata.UpdatedAt = currentTime
		airTranshipmentInfo.Metadata.UpdatedBy = by
		airTranshipmentInfos = append(airTranshipmentInfos, airTranshipmentInfo)
	}

	if len(airTranshipmentInfos) > 0 {

		err = ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocAirTranshipmentTable(ctx)).Save(&airTranshipmentInfos).Error
		if err != nil {
			return nil, err
		}

	}

	doBlInfo := &models.DoBlInfo{}

	if addocs.DoBlInfo.Id == uuid.Nil {

		doBlInfo.Id = uuid.New()
		doBlInfo.ShipmentId = addocs.ShipmentId
		doBlInfo.BlNo = addocs.DoBlInfo.BlNo

		err = ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocDoBlTable(ctx)).Save(&doBlInfo).Error
		if err != nil {
			return nil, err
		}

	}

	return addocs, nil
}

func (b *ArrivalDeliveryDocs) GetArrivalDeliveryDocsById(ctx *context.Context, id uuid.UUID, shipmentId uuid.UUID, HblNumber string) (*models.ArrivalDeliveryDocs, error) {

	addocs := &models.ArrivalDeliveryDocs{}

	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getTable(ctx))
	if id != uuid.Nil {
		tx.Where("id = ?", id)
	}

	if shipmentId != uuid.Nil {
		tx.Where("shipment_id = ?", shipmentId)
	}

	if HblNumber != "" {
		tx.Where("hbl_number = ?", HblNumber)
	}

	err := tx.First(&addocs).Error
	if err != nil {
		return nil, err
	}

	return addocs, nil
}

func (b *ArrivalDeliveryDocs) GetContainerInfoById(ctx *context.Context, id uuid.UUID) ([]*models.ContainerInfo, error) {

	addocscontainerinfo := []*models.ContainerInfo{}

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocContainerTable(ctx)).Where("addoc_id = ?", id).Find(&addocscontainerinfo).Error
	if err != nil {
		return nil, err
	}

	return addocscontainerinfo, nil
}

func (b *ArrivalDeliveryDocs) GetTranshipmentInfoById(ctx *context.Context, id uuid.UUID) ([]*models.TranshipmentInfo, error) {

	addocstranshipmentinfo := []*models.TranshipmentInfo{}

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocTranshipmentTable(ctx)).Where("addoc_id = ?", id).Find(&addocstranshipmentinfo).Error
	if err != nil {
		return nil, err
	}

	return addocstranshipmentinfo, nil
}

func (b *ArrivalDeliveryDocs) GetAirTranshipmentInfo(ctx *context.Context, id uuid.UUID) ([]*models.AirTranshipmentInfo, error) {

	addocsairtranshipmentinfo := []*models.AirTranshipmentInfo{}

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocAirTranshipmentTable(ctx)).Where("addoc_id = ?", id).Find(&addocsairtranshipmentinfo).Error
	if err != nil {
		return nil, err
	}

	return addocsairtranshipmentinfo, nil
}

func (b *ArrivalDeliveryDocs) GetDoBlInfo(ctx *context.Context, id uuid.UUID, blno string) (*models.DoBlInfo, error) {

	addocsdoblinfo := &models.DoBlInfo{}

	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocDoBlTable(ctx))

	if id != uuid.Nil {
		tx.Where("shipment_id = ?", id)
	}

	if blno != "" {
		tx.Where("bl_no = ?", blno)
	}

	err := tx.First(&addocsdoblinfo).Error
	if err != nil {
		return nil, err
	}

	return addocsdoblinfo, nil
}

func (b *ArrivalDeliveryDocs) GetMaxDoBlNumber(ctx *context.Context) int64 {

	var number int64

	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocDoBlTable(ctx)).Select("max(do_no)")
	err := tx.Scan(&number).Error
	if err != nil {
		return 0
	}

	return number
}

func (b *ArrivalDeliveryDocs) DeleteAddocsArrivalDeliveryDocsById(ctx *context.Context, shipmentId uuid.UUID) error {

	var result models.ArrivalDeliveryDocs
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getTable(ctx)).Delete(&result, "shipment_id = ?", shipmentId).Error
	if err != nil {
		ctx.Log.Error("unable to delete arrival delivery docs by shipment id", zap.Error(err))
		return err
	}

	return err
}

func (b *ArrivalDeliveryDocs) DeleteAddocsTranshipmentInfo(ctx *context.Context, addocId uuid.UUID) error {

	var result models.TranshipmentInfo
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocTranshipmentTable(ctx)).Delete(&result, "addoc_id = ?", addocId).Error
	if err != nil {
		ctx.Log.Error("unable to delete arrival delivery docs transhipment info by ad doc id", zap.Error(err))
		return err
	}

	return err
}

func (b *ArrivalDeliveryDocs) DeleteAddocsAirTranshipmentInfo(ctx *context.Context, addocId uuid.UUID) error {

	var result models.AirTranshipmentInfo
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocAirTranshipmentTable(ctx)).Delete(&result, "addoc_id = ?", addocId).Error
	if err != nil {
		ctx.Log.Error("unable to delete arrival delivery docs air transhipment info by ad doc id", zap.Error(err))
		return err
	}

	return err
}

func (b *ArrivalDeliveryDocs) DeleteAddocsContainerInfo(ctx *context.Context, addocId uuid.UUID) error {

	var result models.ContainerInfo
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(b.getAdDocContainerTable(ctx)).Delete(&result, "addoc_id = ?", addocId).Error
	if err != nil {
		ctx.Log.Error("unable to delete arrival delivery docs container info by ad doc id", zap.Error(err))
		return err
	}

	return err
}
