package mastercontainer

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IMasterContainer interface {
	Upsert(ctx *context.Context, m ...*models.MasterContainer) error
	UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.MasterContainer) error
	Update(ctx *context.Context, m *models.MasterContainer) error
	UpdateWithTx(ctx *context.Context, tx *gorm.DB, m *models.MasterContainer) error
	Get(ctx *context.Context, id string) (*models.MasterContainer, error)
	GetAll(ctx *context.Context, ids []string, shipmentId string) ([]*models.MasterContainer, error)
	Delete(ctx *context.Context, id string) error
	DeleteByShipmentId(ctx *context.Context, shipmentId string) error
	DeleteByShipmentIdAndContainerNo(ctx *context.Context, shipmentId, containerNo string) error
}

type MasterContainer struct {
}

func NewMasterContainer() IMasterContainer {
	return &MasterContainer{}
}

func (t *MasterContainer) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "master_containers"
}

func (t *MasterContainer) Upsert(ctx *context.Context, m ...*models.MasterContainer) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *MasterContainer) UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.MasterContainer) error {
	return tx.Table(t.getTable(ctx)).Debug().Save(m).Error
}

func (t *MasterContainer) Update(ctx *context.Context, m *models.MasterContainer) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Debug().Where("id = ?", m.Id).Updates(m).Error
}

func (t *MasterContainer) UpdateWithTx(ctx *context.Context, tx *gorm.DB, m *models.MasterContainer) error {
	return tx.Table(t.getTable(ctx)).Debug().Where("id = ?", m.Id).Updates(m).Error
}

func (t *MasterContainer) Get(ctx *context.Context, id string) (*models.MasterContainer, error) {
	var result models.MasterContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get mastercontainer.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *MasterContainer) Delete(ctx *context.Context, id string) error {
	var result models.MasterContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete mastercontainer.", zap.Error(err))
		return err
	}

	return err
}

func (t *MasterContainer) DeleteByShipmentId(ctx *context.Context, shipmentId string) error {
	var result models.MasterContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "shipment_id = ?", shipmentId).Error
	if err != nil {
		ctx.Log.Error("Unable to delete mastercontainers.", zap.Error(err))
		return err
	}

	return err
}

func (t *MasterContainer) GetAll(ctx *context.Context, ids []string, shipmentId string) ([]*models.MasterContainer, error) {
	var result []*models.MasterContainer

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if len(ids) > 0 {
		tx.Where("id IN ?", ids)
	}

	if shipmentId != "" {
		tx.Where("shipment_id = ?", shipmentId)
	}

	tx.Order("created_at")

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get mastercontainers.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *MasterContainer) DeleteByShipmentIdAndContainerNo(ctx *context.Context, shipmentId, containerNo string) error {
	var result models.MasterContainer
	
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "shipment_id = ? and container_number = ?", shipmentId, containerNo).Error
	if err != nil {
		ctx.Log.Error("Unable to delete mastercontainers.", zap.Error(err))
		return err
	}

	return err
}
