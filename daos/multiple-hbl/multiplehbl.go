package multiplehbl

import (
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IMultipleHbl interface {
	Upsert(ctx *context.Context, m ...*models.MultipleHbl) error
	UpsertMultipleHbl(ctx *context.Context, MultipleHbl *models.MultipleHbl, by uuid.UUID) error
	Get(ctx *context.Context, id string, blNo string, isDeleted *bool, query string) (*models.MultipleHbl, error)
	GetMultipleHbl(ctx *context.Context, blNo string, bookingId uuid.UUID, isDeleted *bool, query string) (*models.MultipleHbl, error)
	GetAll(ctx *context.Context, ids []string, shipmentId string, query string) ([]*models.MultipleHbl, error)
	Delete(ctx *context.Context, id string) error
	GetAllWithoutGeneratedCheck(ctx *context.Context, shipmentId string, isDeleted *bool, query string) ([]*models.MultipleHbl, error)
}

type MultipleHbl struct {
}

func NewMultipleHbl() IMultipleHbl {
	return &MultipleHbl{}
}

func (t *MultipleHbl) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "multiple_hbls"
}

func (t *MultipleHbl) Upsert(ctx *context.Context, m ...*models.MultipleHbl) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *MultipleHbl) UpsertMultipleHbl(ctx *context.Context, MultipleHbl *models.MultipleHbl, by uuid.UUID) error {

	currentTime := time.Now().UTC()

	if MultipleHbl.Id == uuid.Nil {
		MultipleHbl.Id = uuid.New()
		MultipleHbl.CreatedAt = currentTime
		MultipleHbl.CreatedBy = by
	}

	MultipleHbl.UpdatedAt = currentTime
	MultipleHbl.UpdatedBy = by

	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(MultipleHbl).Error
	if err != nil {
		ctx.Log.Error("Failed to upsert MultipleHbl", zap.Error(err))
		return err
	}

	return nil
}

func (t *MultipleHbl) Get(ctx *context.Context, id string, blNo string, isDeleted *bool, query string) (*models.MultipleHbl, error) {
	var result models.MultipleHbl

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if id != "" {
		tx.Where("id = ?", id)
	} else if blNo != "" {
		tx.Where("bl_no = ?", blNo)
	}

	if isDeleted != nil {
		tx.Where("is_deleted = ?", *isDeleted)
	}

	if query != "" {
		tx.Where(query)
	}

	err := tx.First(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get multiplehbl.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (a *MultipleHbl) GetMultipleHbl(ctx *context.Context, blNo string, bookingId uuid.UUID, isDeleted *bool, query string) (*models.MultipleHbl, error) {

	var result models.MultipleHbl

	tx := ctx.DB.WithContext(ctx.Request.Context()).Debug().Table(a.getTable(ctx))

	if blNo != "" {
		tx = tx.Where("bl_no = ?", blNo)
	}

	if bookingId != uuid.Nil {
		tx = tx.Where("shipment_id = ?", bookingId)
	}

	if isDeleted != nil {
		tx = tx.Where("is_deleted = ?", *isDeleted)
	}

	if query != "" {
		tx = tx.Where(query)
	}

	err := tx.First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (t *MultipleHbl) Delete(ctx *context.Context, id string) error {
	var result models.MultipleHbl
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete multiplehbl.", zap.Error(err))
		return err
	}

	return err
}

func (t *MultipleHbl) GetAll(ctx *context.Context, ids []string, shipmentId string, query string) ([]*models.MultipleHbl, error) {
	var result []*models.MultipleHbl

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if len(ids) > 0 {
		tx.Where("id IN ?", ids)
	}

	if shipmentId != "" {
		tx.Where("shipment_id = ?", shipmentId)
	}

	if query != "" {
		tx.Where(query)
	}

	tx.Order("created_at")

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get multiplehbls.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *MultipleHbl) GetAllWithoutGeneratedCheck(ctx *context.Context, shipmentId string, isDeleted *bool, query string) ([]*models.MultipleHbl, error) {

	multipleHbls := []*models.MultipleHbl{}

	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if shipmentId != "" {
		tx.Where("shipment_id = ?", shipmentId)
	}

	if isDeleted != nil {
		tx.Where("is_deleted = ?", *isDeleted)
	}

	if query != "" {
		tx.Where(query)
	}

	tx.Order("created_at")

	err := tx.Find(&multipleHbls).Error
	if err != nil {
		return nil, err
	}

	return multipleHbls, nil
}
