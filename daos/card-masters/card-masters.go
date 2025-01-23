package cardmasters

import (
	"encoding/json"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"

	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"go.uber.org/zap"
)

type ICardMaster interface {
	GetForMilestoneAndTask(ctx *context.Context, milestone, task string) ([]*models.CardMaster, error)
	GetForMilestone(ctx *context.Context, milestone string) ([]*models.CardMaster, error)
	GetForTask(ctx *context.Context, task string) ([]*models.CardMaster, error)
}

type CardMaster struct {
}

func NewCardMaster() ICardMaster {
	return &CardMaster{}
}

func (t *CardMaster) getTable(ctx *context.Context) string {
	return ctx.TenantID + "." + "card_masters"
}

func (t *CardMaster) GetForMilestone(ctx *context.Context, milestone string) ([]*models.CardMaster, error) {
	var dbResults []*models.CardMasterDB
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("milestone = ?", milestone).Find(&dbResults).Error
	if err != nil {
		ctx.Log.Error("Unable to get card masters.", zap.Error(err))
		return nil, err
	}

	var results []*models.CardMaster
	for _, dbResult := range dbResults {
		var info models.CardMasterInfo
		if err := json.Unmarshal(dbResult.Info, &info); err != nil {
			return nil, err
		}

		result := &models.CardMaster{
			Milestone: dbResult.Milestone,
			Task:      dbResult.Task,
			Info:      &info,
		}
		results = append(results, result)
	}

	return results, err
}

func (t *CardMaster) GetForTask(ctx *context.Context, task string) ([]*models.CardMaster, error) {
	var dbResults []*models.CardMasterDB
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("task = ?", task).Find(&dbResults).Error
	if err != nil {
		ctx.Log.Error("Unable to get card masters.", zap.Error(err))
		return nil, err
	}

	var results []*models.CardMaster
	for _, dbResult := range dbResults {
		var info models.CardMasterInfo
		if err := json.Unmarshal(dbResult.Info, &info); err != nil {
			return nil, err
		}

		result := &models.CardMaster{
			Milestone: dbResult.Milestone,
			Task:      dbResult.Task,
			Info:      &info,
		}
		results = append(results, result)
	}

	return results, err
}

func (t *CardMaster) GetForMilestoneAndTask(ctx *context.Context, milestone, task string) ([]*models.CardMaster, error) {
	var dbResults []*models.CardMasterDB
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("milestone = ?", milestone).Where("task = ?", task).Find(&dbResults).Error
	if err != nil {
		ctx.Log.Error("Unable to get card masters.", zap.Error(err))
		return nil, err
	}

	var results []*models.CardMaster
	for _, dbResult := range dbResults {
		var info models.CardMasterInfo
		if err := json.Unmarshal(dbResult.Info, &info); err != nil {
			return nil, err
		}

		result := &models.CardMaster{
			Milestone: dbResult.Milestone,
			Task:      dbResult.Task,
			Info:      &info,
		}
		results = append(results, result)
	}

	return results, err
}
