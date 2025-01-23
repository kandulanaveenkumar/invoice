package shipmentcontact

// import (
// 	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
// 	"bitbucket.org/radarventures/forwarder-shipments/constants"
// 	"bitbucket.org/radarventures/forwarder-shipments/database/models"
// 	"github.com/google/uuid"
// 	"go.uber.org/zap"
// )

// type IShipmentContact interface {
// 	Upsert(ctx *context.Context, m ...*models.ShipmentContact) error
// 	Get(ctx *context.Context, id string) (*models.ShipmentContact, error)
// 	GetAll(ctx *context.Context, ids []string, shipmentId []string) ([]*models.ShipmentContact, error)
// 	Delete(ctx *context.Context, id string) error

// 	GetByShipment(ctx *context.Context, shipmentId uuid.UUID) ([]*models.ShipmentContact, error)
// 	GetShipmentListSearch(ctx *context.Context) ([]*models.ShipmentSearchFilter, error)
// }

// type ShipmentContact struct {
// }

// func NewShipmentContact() IShipmentContact {
// 	return &ShipmentContact{}
// }

// func (t *ShipmentContact) getTable(ctx *context.Context) string {
// 	return ctx.TenantID + "." + "shipment_contacts"
// }

// func (t *ShipmentContact) Upsert(ctx *context.Context, m ...*models.ShipmentContact) error {
// 	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
// }

// func (t *ShipmentContact) Get(ctx *context.Context, id string) (*models.ShipmentContact, error) {
// 	var result models.ShipmentContact
// 	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
// 	if err != nil {
// 		ctx.Log.Error("Unable to get shipmentcontact.", zap.Error(err))
// 		return nil, err
// 	}

// 	return &result, err
// }

// func (t *ShipmentContact) Delete(ctx *context.Context, id string) error {
// 	var result models.ShipmentContact
// 	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
// 	if err != nil {
// 		ctx.Log.Error("Unable to delete shipmentcontact.", zap.Error(err))
// 		return err
// 	}

// 	return err
// }

// func (t *ShipmentContact) GetAll(ctx *context.Context, ids []string, shipmentId []string) ([]*models.ShipmentContact, error) {
// 	var result []*models.ShipmentContact

// 	tx := ctx.DB.Debug().Debug().Table(t.getTable(ctx))

// 	if len(ids) > 0 {
// 		tx.Where("id IN ?", ids)
// 	}

// 	if len(shipmentId) > 0 {
// 		tx.Where("shipment_id IN ?", shipmentId)
// 	}

// 	err := tx.Find(&result).Error
// 	if err != nil {
// 		ctx.Log.Error("Unable to get shipment contacts", zap.Error(err))
// 		return nil, err
// 	}

// 	return result, err
// }

// func (t *ShipmentContact) GetShipmentListSearch(ctx *context.Context) ([]*models.ShipmentSearchFilter, error) {
// 	var shipmentSearch []*models.ShipmentSearchFilter
// 	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("id", "name", "shipment_contacts.type", "email", "mobile", "shipment_contacts.created_at").
// 		Joins("JOIN "+ctx.TenantID+".shipments ON shipments.id = shipment_contacts.shipment_id").
// 		Where("is_deleted = false AND ? = ANY(ARRAY[region_id,origin_region_id,dest_region_id])", ctx.Account.RegionID).
// 		Where("shipment_contacts.type IN (?)", []string{constants.ShipmentContactShipper, constants.ShipmentContactConsignee}).
// 		Where("name ilike ?", ctx.Query("q")+"%").
// 		Scan(&shipmentSearch).Error
// 	if err != nil {
// 		ctx.Log.Error("Unable to get shipmentcontainer.", zap.Error(err))
// 		return nil, err
// 	}
// 	return shipmentSearch, nil
// }

// func (t *ShipmentContact) GetByShipment(ctx *context.Context, shipmentId uuid.UUID) ([]*models.ShipmentContact, error) {
// 	var result []*models.ShipmentContact
// 	err := ctx.DB.Table(t.getTable(ctx)).Where("shipment_id = ?", shipmentId).Find(&result).Error
// 	if err != nil {
// 		ctx.Log.Error("Unable to get shipmentcontacts.", zap.Error(err))
// 		return nil, err
// 	}

// 	return result, err
// }
