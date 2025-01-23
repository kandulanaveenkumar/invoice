package execassign

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type IExecAssign interface {
	Upsert(ctx *context.Context, m ...*models.ExecAssign) error
	Get(ctx *context.Context, execType, executiveId string) (*models.ExecAssign, error)
	Delete(ctx *context.Context, execType, executiveId string) error
}

type ExecAssign struct {
}

func NewExecAssign() IExecAssign {
	return &ExecAssign{}
}

func (t *ExecAssign) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "exec_assign"
}

func (t *ExecAssign) Upsert(ctx *context.Context, m ...*models.ExecAssign) error {
	uniqueConstraint := "exec_assign_pkey"
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Clauses(clause.OnConflict{OnConstraint: uniqueConstraint, UpdateAll: true}).Save(m).Error
}

func (t *ExecAssign) Get(ctx *context.Context, execType, executiveId string) (*models.ExecAssign, error) {
	var result models.ExecAssign
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "exec_type = ? AND executive_id = ?", execType, executiveId).Error
	if err != nil {
		ctx.Log.Error("Unable to get exec_assign.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *ExecAssign) Delete(ctx *context.Context, execType, executiveId string) error {
	var result models.ExecAssign
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "exec_type = ? AND executive_id = ?", execType, executiveId).Error
	if err != nil {
		ctx.Log.Error("Unable to delete exec_assign.", zap.Error(err))
		return err
	}

	return err
}
