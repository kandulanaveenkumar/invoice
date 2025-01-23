package workflow

import (
	"fmt"
	"time"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type IFlowInstances interface {
	Upsert(ctx *context.Context, m *models.FlowInstances) error
	Get(ctx *context.Context, id string) (*models.FlowInstances, error)
	Delete(ctx *context.Context, id string) error
	GetForInstance(ctx *context.Context, instanceID, flowType string) ([]*models.FlowInstances, error)
	GetCardLookUpDetails(ctx *context.Context, req dtos.WorkflowInstanceReq) ([]dtos.CardLabelLists, error)
	GetDistinctColumnDetailsRfqData(ctx *context.Context, req dtos.WorkflowInstanceReq) ([]dtos.CardLabelLists, error)
	GetDistinctColumnDetailsShipmentData(ctx *context.Context, req dtos.WorkflowInstanceReq) ([]dtos.CardLabelLists, error)
	GetCardInstancesWithIds(ctx *context.Context, ids []string) ([]models.FlowInstances, error)
	ReassignAllBookingAssignedCard(ctx *context.Context, req *dtos.ReExecCard, instance_ids []string, aid string) ([]models.FlowInstances, error)
	ReassignCard(ctx *context.Context, req *dtos.ReExecCard) (*models.FlowInstances, error)
	EscalateCard(ctx *context.Context, filter *dtos.ReExecCard, EscalatedByList pq.StringArray, id string) error
	GetExecutiveRfqBookingsCardsCount(ctx *context.Context, InstanceIds []string, ExecIds []string, InstanceType string) ([]*models.ExecutiveCardCount, error)
	GetNonEscalatedCards(ctx *context.Context, filter *models.FlowInstances, statuslist []string, nonesc bool) ([]models.FlowInstances, error)
	GetCardswithNoBooks(ctx *context.Context, rep_ids []string, cardFilter *dtos.GECardFilter) (models.CountCards, error)
	GetLiveEscaltedCardsOrgtree(ctx *context.Context) (map[string]dtos.TaskCounts, error)
	GetCardCounts(ctx *context.Context) (map[string]dtos.TaskCounts, error)
	DeleteCardInstances(ctx *context.Context, id string) error
}

type FlowInstances struct {
}

func NewFlowInstances() IFlowInstances {
	return &FlowInstances{}
}

func (t *FlowInstances) getTable(ctx *context.Context) string {
	return ctx.TenantID + ".flow_instances"
}

func (t *FlowInstances) getShipmentTable(ctx *context.Context) string {
	return ctx.TenantID + ".shipments"

}
func (t *FlowInstances) getRfqTable(ctx *context.Context) string {
	return ctx.TenantID + ".rfqs"

}
func (t *FlowInstances) Upsert(ctx *context.Context, m *models.FlowInstances) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *FlowInstances) Get(ctx *context.Context, id string) (*models.FlowInstances, error) {
	var result models.FlowInstances
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *FlowInstances) Delete(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Unscoped().Delete(&models.FlowInstances{
		Id: id,
	}).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow field.", zap.Error(err))
		return err
	}

	return err
}

func (t *FlowInstances) GetForInstance(ctx *context.Context, instanceID, flowType string) ([]*models.FlowInstances, error) {
	var result []*models.FlowInstances
	tx := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("instance_id = ?", instanceID)
	if flowType != "" {
		tx = tx.Where("type = ?", flowType)
	}

	err := tx.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow fields.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *FlowInstances) GetCardLookUpDetails(ctx *context.Context, req dtos.WorkflowInstanceReq) ([]dtos.CardLabelLists, error) {

	var results []dtos.CardLabelLists

	currentTime := time.Now().UTC()

	query := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Distinct(fmt.Sprintf(" %s.%s AS label , %s.%s AS value,", t.getTable(ctx), req.CardLabelLists.Label, t.getTable(ctx), req.CardLabelLists.Value)+"'"+req.CardLabelLists.Type+"' AS type ").
		Where("card_status in ?", req.StatusList).
		Where("visible_time < ?", currentTime).
		Where(" assigned_to in (?) ", req.AssignedTo).
		Where("card_status IN (?)", req.StatusList).
		Where(req.CardLabelLists.Label+" ILIKE (?)", "%"+req.Qeurry+"%")
	err := query.Find(&results).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow instance look up", zap.Error(err))
	}

	return results, nil

}

func (t *FlowInstances) GetDistinctColumnDetailsRfqData(ctx *context.Context, req dtos.WorkflowInstanceReq) ([]dtos.CardLabelLists, error) {

	var results []dtos.CardLabelLists

	currentTime := time.Now().UTC()

	query := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Distinct(fmt.Sprintf(" %s.%s AS label , %s.%s AS value,", t.getRfqTable(ctx), req.CardLabelLists.Label, t.getRfqTable(ctx), req.CardLabelLists.Value)+"'"+req.CardLabelLists.Type+"' AS type ").
		Joins(fmt.Sprintf(" JOIN %s on %s.id::TEXT=%s.instance_id", t.getRfqTable(ctx), t.getRfqTable(ctx), t.getTable(ctx))).
		Where("card_status in ?", req.StatusList).
		Where("visible_time < ?", currentTime).
		Where(" assigned_to in (?) ", req.AssignedTo)
	if req.Qeurry != "" {
		query.Where(fmt.Sprintf(" %s.%s ILIKE (?)", t.getTable(ctx), req.CardLabelLists.Label), "%"+req.Qeurry+"%")
	}
	err := query.Find(&results).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow instance look up", zap.Error(err))
	}
	return results, nil

}

func (t *FlowInstances) GetDistinctColumnDetailsShipmentData(ctx *context.Context, req dtos.WorkflowInstanceReq) ([]dtos.CardLabelLists, error) {

	var results []dtos.CardLabelLists

	currentTime := time.Now().UTC()

	query := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Distinct(fmt.Sprintf(" %s.%s AS label , %s.%s AS value,", t.getShipmentTable(ctx), req.CardLabelLists.Label, t.getShipmentTable(ctx), req.CardLabelLists.Value)+"'"+req.CardLabelLists.Type+"' AS type ").
		Joins(fmt.Sprintf(" JOIN %s on %s.id::TEXT=%s.instance_id::TEXT", t.getShipmentTable(ctx), t.getShipmentTable(ctx), t.getTable(ctx))).
		Where("card_status in ?", req.StatusList).
		Where("visible_time < ?", currentTime).
		Where(" assigned_to in (?) ", req.AssignedTo)
	if req.Qeurry != "" {
		query.Where(fmt.Sprintf(" %s.%s ILIKE (?)", t.getTable(ctx), req.CardLabelLists.Label), "%"+req.Qeurry+"%")
	}
	err := query.Find(&results).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow instance look up", zap.Error(err))
	}
	return results, nil

}

func (t *FlowInstances) GetCardInstancesWithIds(ctx *context.Context, ids []string) ([]models.FlowInstances, error) {
	var result []models.FlowInstances
	query := ctx.DB.Table(t.getTable(ctx)).Where("id in ?", ids)
	err := query.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get workflow instance with ids", zap.Error(err))
	}
	return result, err
}
func (t *FlowInstances) ReassignAllBookingAssignedCard(ctx *context.Context, req *dtos.ReExecCard, instance_ids []string, aid string) ([]models.FlowInstances, error) {
	statusList := []string{globals.StatusBreached, globals.StatusCreated, globals.StatusWarning}
	newValues := map[string]interface{}{
		"assigned_to":            req.AssignTo,
		"assigned_to_name":       req.AssignToName,
		"escalated":              false,
		"escalated_to":           "",
		"escalated_by_id":        pq.StringArray{},
		"escalated_by":           "",
		"escalation_reason":      "",
		"escalation_description": "",
		"updated_at":             time.Now().UTC(),
	}
	err := ctx.DB.Debug().Table(t.getTable(ctx)).Where("status in (?)", statusList).Where("instance_id in (?) and assigned_to = ?", instance_ids, aid).Updates(newValues).Error
	if err != nil {
		ctx.Log.Error("Error while reassigning card in DB", zap.Error(err))
	}
	var updatedCards []models.FlowInstances
	err = ctx.DB.Debug().Table(t.getTable(ctx)).
		Where("status IN (?)", statusList).
		Where("instance_id IN (?) AND assigned_to = ?", instance_ids, aid).
		Find(&updatedCards).
		Error
	if err != nil {
		ctx.Log.Error("Error while retrieving updated cards from DB", zap.Error(err))
		return nil, err
	}
	return updatedCards, nil
}
func (t *FlowInstances) ReassignCard(ctx *context.Context, req *dtos.ReExecCard) (*models.FlowInstances, error) {
	flow_instance_ids := req.CardIds
	statusList := []string{globals.StatusBreached, globals.StatusCreated, globals.StatusWarning}
	newValues := map[string]interface{}{
		"assigned_to":            req.AssignTo,
		"assigned_to_name":       req.AssignToName,
		"escalated":              false,
		"escalated_to":           "",
		"escalated_by_id":        pq.StringArray{},
		"escalated_by":           "",
		"escalation_reason":      "",
		"escalation_description": "",
		"updated_at":             time.Now().UTC(),
	}
	err := ctx.DB.Debug().Table(t.getTable(ctx)).Where("status in (?)", statusList).Where("id in (?)", flow_instance_ids).Updates(newValues).Error
	if err != nil {
		ctx.Log.Error("Error while reassigning card in DB", zap.Error(err))
		return nil, err
	}
	var updatedCard *models.FlowInstances
	err = ctx.DB.Debug().Table(t.getTable(ctx)).
		Where("status IN (?)", statusList).
		Where("id in (?)", flow_instance_ids).
		Find(&updatedCard).
		Error
	if err != nil {
		ctx.Log.Error("Error while retrieving updated cards from DB", zap.Error(err))
		return nil, err
	}
	return updatedCard, nil
}
func (t *FlowInstances) EscalateCard(ctx *context.Context, filter *dtos.ReExecCard, EscalatedByList pq.StringArray, id string) error {
	newValues := map[string]interface{}{
		"escalated":              true,
		"escalated_to":           filter.EscalateTo,
		"escalated_by_id":        EscalatedByList,
		"escalated_by":           filter.EscalatedByName,
		"escalation_reason":      filter.EscalationReason,
		"escalation_description": filter.EscalationRemarks,
		"updated_at":             time.Now().UTC(),
	}
	err := ctx.DB.Debug().Table(t.getTable(ctx)).Where("status not in ('Completed','Delete') and id = ?", id).Updates(newValues).Error
	if err != nil {
		ctx.Log.Error("Error while escalating card in DB", zap.Error(err))
	}
	return nil
}

func (t *FlowInstances) GetExecutiveRfqBookingsCardsCount(ctx *context.Context, InstanceIds []string, ExecIds []string, InstanceType string) ([]*models.ExecutiveCardCount, error) {
	var cardsCount []*models.ExecutiveCardCount
	query := ctx.DB.Debug().Raw(`select execs.executive_id,execs.id, count(c.id) from 
	unnest(?::text[],?::text[]) as execs(id,executive_id) left join `+t.getTable(ctx)+` c
	on (c.instance_id  = execs.id and c.assigned_to = execs.executive_id) and c.instance_type=`+InstanceType+` 
	group by execs.executive_id, execs.id`, pq.Array(InstanceIds), pq.Array(ExecIds))
	err := query.Find(&cardsCount).Error
	if err != nil {
		ctx.Log.Error("Error while GetExecutiveCardsCount", zap.Error(err))
		return nil, err
	}
	return cardsCount, nil
}
func (t *FlowInstances) GetNonEscalatedCards(ctx *context.Context, filter *models.FlowInstances, statuslist []string, nonesc bool) ([]models.FlowInstances, error) {
	var results []models.FlowInstances
	query := ctx.DB.Debug().Table(t.getTable(ctx)).Where(&filter)
	if nonesc {
		query = query.Where("escalated = false")
	}
	if len(statuslist) > 0 {
		query = query.Where("status in (?)", statuslist)
	}
	err := query.Find(&results).Error
	if err != nil {
		ctx.Log.Error("Error while FilterCards", zap.Any("filter", filter), zap.Error(err))
		return nil, err
	}
	return results, nil
}

func (t *FlowInstances) GetCardswithNoBooks(ctx *context.Context, rep_ids []string, cardFilter *dtos.GECardFilter) (models.CountCards, error) {
	var results models.CountCards
	currentTime := time.Now().UTC()
	statusList := []string{globals.StatusBreached, globals.StatusCreated, globals.StatusWarning}
	query := ctx.DB.Debug().Table(t.getTable(ctx)).
		Select("assigned_to,assigned_to_name ,status, COUNT(DISTINCT(id)) as task_count").
		Where("assigned_to in ?", rep_ids).
		Where("visible_time < ?", currentTime).
		Where("status in (?)", statusList).
		Where("instance_id = ''").
		Group("assigned_to, assigned_to_name, status")
	err := query.Find(&results).Error
	if err != nil {
		ctx.Log.Error("Error while getting bulk  --DB", zap.Error(err))
	}
	return results, err
}

func (t *FlowInstances) GetLiveEscaltedCardsOrgtree(ctx *context.Context) (map[string]dtos.TaskCounts, error) {
	rows, err := ctx.DB.Debug().Raw(`SELECT
    unnested_uuid AS escalated_by_id,
    COUNT(id) AS uuid_count
FROM
    (
        SELECT
            id,assigned_to,
            unnest(escalated_by_id) AS unnested_uuid
        FROM
            ` + t.getTable(ctx) + `
        WHERE
            escalated = true
            AND status NOT IN ('Delete', 'Completed')
    ) AS subquery
    where subquery.assigned_to =subquery.unnested_uuid
GROUP BY
    unnested_uuid
ORDER BY
    uuid_count DESC`).Rows()
	if err != nil {
		ctx.Log.Error("error while getting GetLiveEscaltedCards column details", zap.Error(err))
		return nil, err
	}
	TaskCounts := make(map[string]dtos.TaskCounts)
	for rows.Next() {
		var id string
		var counts dtos.TaskCounts
		err = rows.Scan(&id, &counts.EscalatedTasks)
		if err != nil {
			return nil, err
		}
		TaskCounts[id] = counts
	}
	return TaskCounts, nil
}

func (t *FlowInstances) GetCardCounts(ctx *context.Context) (map[string]dtos.TaskCounts, error) {
	rows, err := ctx.DB.Debug().Raw(`SELECT assigned_to,
	Sum(CASE
		  WHEN ( status IN( 'Created', 'Warning') ) THEN 1
		  ELSE 0
		END) AS pending,
	Sum(CASE
		  WHEN ( status = 'Breached' ) THEN 1
		  ELSE 0
		END) AS dealyed,
	Sum(CASE
		  WHEN ( status = 'Warning' ) THEN 1
		  ELSE 0
		END) AS warning
FROM  ` + t.getTable(ctx) + `
WHERE  visible_time < Now()
	AND escalated = false
	AND instance_id != ''
GROUP  BY assigned_to `).Rows()
	if err != nil {
		ctx.Log.Error("Error while Getting team cards from DB", zap.Error(err))
	}
	defer rows.Close()
	taskCounts := make(map[string]dtos.TaskCounts)
	for rows.Next() {
		var id string
		var counts dtos.TaskCounts
		err = rows.Scan(&id, &counts.PendingTasks, &counts.BreachedTasks, &counts.ExpiringTasks)
		if err != nil {
			return nil, err
		}
		taskCounts[id] = counts
	}
	return taskCounts, nil
}

func (t *FlowInstances) DeleteCardInstances(ctx *context.Context, id string) error {
	err := ctx.DB.Debug().Raw(`Update` + t.getTable(ctx) + `set active=false where instance_type = 'shipment' and instance_id = ` + id + ``).Error
	if err != nil {
		ctx.Log.Error("unable to delte cards instances", zap.Any("", err))
		return err
	}

	return nil
}
