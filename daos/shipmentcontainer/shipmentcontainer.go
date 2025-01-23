package shipmentcontainer

import (
	"fmt"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IShipmentContainer interface {
	Upsert(ctx *context.Context, m ...*models.ShipmentContainer) error
	UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.ShipmentContainer) error
	UpdateWithTx(ctx *context.Context, tx *gorm.DB, m *models.ShipmentContainer) error
	Update(ctx *context.Context, m *models.ShipmentContainer) error

	Get(ctx *context.Context, id string) (*models.ShipmentContainer, error)
	GetAll(ctx *context.Context, shipmentIds []string) ([]*models.ShipmentContainer, error)
	Delete(ctx *context.Context, id string) error
	DeleteByShipmentId(ctx *context.Context, shipmentId string) error

	GetByShipment(ctx *context.Context, shipmentIds []string) ([]*models.ShipmentContainer, error)
	GetShipmentListSearch(ctx *context.Context) ([]*models.ShipmentSearchFilter, error)

	GetByShipmentAndFlowInstance(ctx *context.Context, shipmentId, instanceID string) ([]*models.ShipmentContainer, error)
	UpdateFlowInstanceID(ctx *context.Context, id, instanceID string) error
	GetCountByShipmentId(ctx *context.Context, shipmentId string) int64
	DeleteContainer(ctx *context.Context, m *models.ShipmentContainer) error
}

type ShipmentContainer struct {
}

func NewShipmentContainer() IShipmentContainer {
	return &ShipmentContainer{}
}

func (t *ShipmentContainer) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "shipment_containers"
}

func (t *ShipmentContainer) Upsert(ctx *context.Context, m ...*models.ShipmentContainer) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Debug().Error
}

func (t *ShipmentContainer) DeleteContainer(ctx *context.Context, m *models.ShipmentContainer) error {
	return ctx.DB.WithContext(ctx.Request.Context()).
		Table(t.getTable(ctx)).
		Where("id = ?", m.Id.String()).
		Updates(map[string]interface{}{
			"is_deleted": m.IsDeleted,
			"is_active":  false,
		}).Error
}

func (t *ShipmentContainer) UpsertWithTx(ctx *context.Context, tx *gorm.DB, m ...*models.ShipmentContainer) error {
	return tx.Table(t.getTable(ctx)).Debug().Save(m).Error
}

func (t *ShipmentContainer) Get(ctx *context.Context, id string) (*models.ShipmentContainer, error) {
	var result models.ShipmentContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentcontainer.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *ShipmentContainer) Delete(ctx *context.Context, id string) error {
	var result models.ShipmentContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete shipmentcontainer.", zap.Error(err))
		return err
	}

	return err
}

func (t *ShipmentContainer) DeleteByShipmentId(ctx *context.Context, shipmentId string) error {
	var result models.ShipmentContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "shipment_id = ?", shipmentId).Error
	if err != nil {
		ctx.Log.Error("Unable to delete shipmentcontainers.", zap.Error(err))
		return err
	}

	return err
}

func (t *ShipmentContainer) GetAll(ctx *context.Context, shipmentIds []string) ([]*models.ShipmentContainer, error) {
	var result []*models.ShipmentContainer
	if len(shipmentIds) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get shipmentcontainers.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).
		Table(t.getTable(ctx)).
		Where("shipment_id IN (?) AND is_deleted = false", shipmentIds).
		Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentcontainers.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *ShipmentContainer) GetByShipment(ctx *context.Context, shipmentIds []string) ([]*models.ShipmentContainer, error) {
	var result []*models.ShipmentContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result, "shipment_id IN (?) AND is_deleted = false", shipmentIds).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentcontainer.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *ShipmentContainer) GetShipmentListSearch(ctx *context.Context) ([]*models.ShipmentSearchFilter, error) {
	var shipmentSearch []*models.ShipmentSearchFilter
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("shipments.id, no as code, shipment_containers.created_at").
		Joins("JOIN "+ctx.TenantID+".shipments ON shipments.id = shipment_containers.shipment_id ").
		Where("shipments.is_deleted = false AND ? = ANY(ARRAY[region_id,origin_region_id,dest_region_id])", ctx.Account.RegionID).
		Where("shipment_containers.is_deleted = false AND no ilike ?", ctx.Query("q")+"%").
		Scan(&shipmentSearch).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipmentcontainer.", zap.Error(err))
		return nil, err
	}
	return shipmentSearch, nil
}

func (t *ShipmentContainer) GetByShipmentAndFlowInstance(ctx *context.Context, shipmentId, instanceID string) ([]*models.ShipmentContainer, error) {
	var result []*models.ShipmentContainer
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result, "shipment_id = ? AND flow_instance_id = ?", shipmentId, instanceID).Error
	if err != nil {
		ctx.Log.Error("Unable to get shipment containers for flow instance.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *ShipmentContainer) UpdateFlowInstanceID(ctx *context.Context, id, instanceID string) error {
	err := ctx.DB.Exec(fmt.Sprintf(`UPDATE %s set flow_instance_id = ? WHERE id = ? `, t.getTable(ctx)), instanceID, id).Error
	if err != nil {
		ctx.Log.Error("Unable to update shipment container's flow instance.", zap.Error(err))
		return err
	}

	return err
}

func (t *ShipmentContainer) GetCountByShipmentId(ctx *context.Context, shipmentId string) int64 {

	var count int64
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("shipment_id = ? AND is_deleted = false", shipmentId).Count(&count).Error
	if err != nil {
		ctx.Log.Error("unable to get shipment containers count", zap.Error(err))
		return 0
	}

	return count
}

func (t *ShipmentContainer) UpdateWithTx(ctx *context.Context, tx *gorm.DB, m *models.ShipmentContainer) error {
	return tx.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Debug().Where("id = ?", m.Id).Updates(m).Error
}

func (t *ShipmentContainer) Update(ctx *context.Context, m *models.ShipmentContainer) error {
	return ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Debug().Where("id = ?", m.Id).Updates(m).Error
}
