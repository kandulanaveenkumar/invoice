package card

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"go.uber.org/zap"
)

func (t *Card) GetPendingTasksCountForExec(ctx *context.Context, execID, name string) (int, error) {
	var count int
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getRfqsTable(ctx)+" as rfqs").
		Select("count(c.id)").
		Joins("JOIN cards c ON (rfqs.id)::TEXT = c.instance_id AND c.name = ? AND c.assigned_to = ?", name, execID).
		Joins("LEFT JOIN cards c2 ON (rfqs.id)::TEXT = c2.instance_id AND c2.name = ? AND c2.assigned_to != ?", name, execID).
		Where("rfqs.status ILIKE ?", "%"+constants.RfqStatusBuyTBA).
		Where("(c.status != ? AND c.status != ? AND c.status != ?) OR (c2.status != ? AND c2.status != ? AND c.status != ?)",
			constants.CardStatusBreached, constants.CardStatusCompleted, constants.CardStatusDeleted,
			constants.CardStatusBreached, constants.CardStatusCompleted, constants.CardStatusDeleted).
		Scan(&count).Error

	if err != nil {
		ctx.Log.Error("Unable to get pending tasks count for executive.", zap.Error(err))
		return 0, err
	}

	return count, err
}

func (t *Card) GetExpiredTasksCountForExec(ctx *context.Context, execID, name string) (int, error) {
	var result int
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Joins("JOIN rfqs on ((rfqs.id)::TEXT = cards.instance_id)").
		Select("count(cards.id)").Where("cards.assigned_to = ?", execID).
		Where("rfqs.status ilike ?", "%"+constants.RfqStatusBuyTBA).
		Where("cards.name = ?", name).Where("cards.status = ?", constants.CardStatusBreached).Scan(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get card.", zap.Error(err))
		return 0, err
	}

	return result, err
}
