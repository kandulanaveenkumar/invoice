package document

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IDocument interface {
	Upsert(ctx *context.Context, m ...*models.Document) error
	Get(ctx *context.Context, id uuid.UUID) (*models.Document, error)
	GetForShipment(ctx *context.Context, shipmentId string, documentIds []string, name, owner []string, q, regionId string) ([]*models.Document, error)
	GetWithFilters(ctx *context.Context, documentIds, instanceIds []uuid.UUID, name, owner []string, nameLike string) ([]*models.Document, error)
	DeleteByInstanceId(ctx *context.Context, shipmentId string) error
	DeleteByFlowInstanceId(ctx *context.Context, instanceId string, flowInstanceId string) error
}

type Document struct {
}

func NewDocument() IDocument {
	return &Document{}
}

func (t *Document) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "documents"
}

func (t *Document) Upsert(ctx *context.Context, m ...*models.Document) error {
	return ctx.DB.Table(t.getTable(ctx)).Save(m).Error
}

func (t *Document) Get(ctx *context.Context, id uuid.UUID) (*models.Document, error) {
	var result models.Document
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get documents.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *Document) GetWithFilters(ctx *context.Context, documentIds, instanceIds []uuid.UUID, name, owner []string, nameLike string) ([]*models.Document, error) {

	tx := ctx.DB.Debug().Table(t.getTable(ctx))

	if len(documentIds) > 0 {
		tx.Where("document_id IN (?)", documentIds)
	}

	if len(instanceIds) > 0 {
		tx.Where("instance_id IN (?)", instanceIds)
	}

	if len(name) > 0 {
		tx.Where("name IN (?)", name)
	}

	if len(owner) > 0 {
		tx.Where("owner IN (?)", owner)
	}

	if nameLike != "" {
		tx.Where("name ilike ", nameLike+"%")
	}

	tx.Order("created_at")

	var result []*models.Document

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get documents.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Document) GetForShipment(ctx *context.Context, instanceId string, documentIds []string, name, owner []string, q, regionId string) ([]*models.Document, error) {

	tx := ctx.DB.Debug().Table(t.getTable(ctx))

	if instanceId != "" {
		tx.Where("instance_id = ?", instanceId)
	}

	if len(documentIds) > 0 {
		tx.Where("document_id IN (?)", documentIds)
	}

	if len(name) > 0 {
		tx.Where("name IN (?)", name)
	}

	if len(owner) > 0 {
		tx.Where("owner IN (?)", owner)
	}

	if q != "" {
		tx.Where("metadata->'invoice_no'", q+"%")
		tx.Where("metadata->'invoice_id'", q+"%")
	}

	if regionId != "" {
		tx.Where("region_id IN (?)", regionId)
	}

	tx.Order("created_at")

	var result []*models.Document

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get documents.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Document) DeleteByInstanceId(ctx *context.Context, instanceId string) error {
	var result models.Document
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "instance_id = ?", instanceId).Error
	if err != nil {
		ctx.Log.Error("Unable to delete documents.", zap.Error(err))
		return err
	}

	return err
}

func (t *Document) DeleteByFlowInstanceId(ctx *context.Context, instanceId string, flowInstanceId string) error {
	var result models.Document
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "instance_id = ? AND flow_instance_id = ?", instanceId, flowInstanceId).Error

	if err != nil {
		ctx.Log.Error("Unable to delete documents.", zap.Error(err))
		return err
	}

	return err
}
