package marksanddescription

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type IMarksAndDescription interface {
	Upsert(ctx *context.Context, m ...*models.MarksAndDescription) error
	Get(ctx *context.Context, id string) (*models.MarksAndDescription, error)
	GetAll(ctx *context.Context, ids []string, query string) ([]*models.MarksAndDescription, error)
	Delete(ctx *context.Context, id string) error

	GetForHBL(ctx *context.Context, bl_no string) ([]*models.MarksAndDescription, error)
}

type MarksAndDescription struct {
}

func NewMarksAndDescription() IMarksAndDescription {
	return &MarksAndDescription{}
}

func (t *MarksAndDescription) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "marks_and_descriptions"
}

func (t *MarksAndDescription) Upsert(ctx *context.Context, m ...*models.MarksAndDescription) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *MarksAndDescription) Get(ctx *context.Context, id string) (*models.MarksAndDescription, error) {
	var result models.MarksAndDescription
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get marksanddescription.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *MarksAndDescription) Delete(ctx *context.Context, id string) error {
	var result models.MarksAndDescription
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete marksanddescription.", zap.Error(err))
		return err
	}

	return err
}

func (t *MarksAndDescription) GetAll(ctx *context.Context, ids []string, query string) ([]*models.MarksAndDescription, error) {
	var result []*models.MarksAndDescription

	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx))

	if len(ids) > 0 {
		tx.Where("id IN ?", ids)
	}

	if query != "" {
		tx.Where(query)
	}

	tx.Order("created_at")

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get marksanddescriptions.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *MarksAndDescription) GetForHBL(ctx *context.Context, bl_no string) ([]*models.MarksAndDescription, error) {
	var result []*models.MarksAndDescription
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result, "bl_no = ?", bl_no).Error
	if err != nil {
		ctx.Log.Error("Unable to get marksanddescription.", zap.Error(err))
		return nil, err
	}

	return result, err
}
