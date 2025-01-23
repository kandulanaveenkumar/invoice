package card

import (
	"errors"
	"fmt"
	"strings"
	"time"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type ICard interface {
	Upsert(ctx *context.Context, m ...*models.Card) error
	Get(ctx *context.Context, id string) (*models.Card, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.Card, error)
	Delete(ctx *context.Context, id string) error
	GetCardsWithFilter(ctx *context.Context, filter *models.Card, statuslist []string) ([]*models.Card, error)
	GetExecCards(ctx *context.Context, req *models.Card, cardfilter *dtos.GECardFilter) ([]*models.Card, error)
	GetAssignedTo(ctx *context.Context, filter *models.Card, statusList []string) (*models.Card, error)
	GetShipmentsWithFilter(ctx *context.Context, isShipment bool, rep_ids []string, shipmentFilter *dtos.GEShipmentFilter, cardFilter *dtos.GECardFilter) (models.CountCards, error)
	GetCardswithNoBooks(ctx *context.Context, rep_ids []string, cardFilter *dtos.GECardFilter) (models.CountCards, error)
	GetLiveEscaltedCardsOrgtree(ctx *context.Context) (map[string]dtos.TaskCounts, error)
	GetCardCounts(ctx *context.Context) (map[string]dtos.TaskCounts, error)
	ReassignAllShipmentAssignedCard(ctx *context.Context, req *dtos.ReExecCard, instance_ids []string, aid string) ([]models.Card, error)
	ReassignCard(ctx *context.Context, req *dtos.ReExecCard) (*models.Card, error)
	EscalateCard(ctx *context.Context, filter *dtos.ReExecCard, EscalatedByList pq.StringArray, id string) error
	GetNonEscalatedCards(ctx *context.Context, filter *models.Card, statuslist []string, nonesc bool) ([]models.Card, error)
	GetDistinctColumnDetails(ctx *context.Context, label, value, cardType, querry string, ids, cardStatus []string) ([]dtos.CardLabelLists, error)
	GetDistinctColumnDetailsInstanceData(ctx *context.Context, val dtos.CardLabelLists, filter *dtos.FilterCardLabel, ids, cardStatus []string) ([]dtos.CardLabelLists, error)
	GetCardsUsingMulti(ctx *context.Context, cardIds []string) ([]models.Card, error)
	GetAllPendingCards(ctx *context.Context) ([]models.Card, error)
	UpdateStatus(ctx *context.Context, id, status string) error
	GetCardsFiltered(ctx *context.Context, req *models.Card) ([]*models.Card, error)
	GetBulkCardsByCompanyId(ctx *context.Context, cids []string) ([]*models.Card, error)
	UpsertCardById(ctx *context.Context, m *models.Card, action string) error
	GetPendingTasksCountForExec(ctx *context.Context, execID, name string) (int, error)
	GetExpiredTasksCountForExec(ctx *context.Context, execID, name string) (int, error)
	GetBulkCardsFiltered(ctx *context.Context, filters *models.BulkCardsFilters) ([]*models.Card, error)
	DeleteBookingRequestByID(ctx *context.Context, id string, updatedAt time.Time) (bool, error)
	BulkUpsert(ctx *context.Context, cardIds []string) error
}

type Card struct {
}

func NewCard() ICard {
	return &Card{}
}

func (t *Card) getTable(ctx *context.Context) string {
	if ctx.TenantID == "" {
		ctx.TenantID = "public"
	}
	return ctx.TenantID + "." + "cards"
}

func (t *Card) getRfqsTable(ctx *context.Context) string {
	if ctx.TenantID == "" {
		ctx.TenantID = "public"
	}
	return ctx.TenantID + ".rfqs"
}

func (t *Card) getShipmentsTable(ctx *context.Context) string {
	return ctx.TenantID + ".shipments"
}

func (t *Card) Upsert(ctx *context.Context, m ...*models.Card) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *Card) Get(ctx *context.Context, id string) (*models.Card, error) {
	var result models.Card
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get card.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *Card) Delete(ctx *context.Context, id string) error {
	var result models.Card
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete card.", zap.Error(err))
		return err
	}

	return err
}

func (t *Card) GetAll(ctx *context.Context, ids []string) ([]*models.Card, error) {
	var result []*models.Card
	if len(ids) == 0 {
		err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get cards.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get cards.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Card) GetCardsWithFilter(ctx *context.Context, filter *models.Card, statuslist []string) ([]*models.Card, error) {
	var results []*models.Card

	nameList := []string{}
	name := filter.Name

	if name != "" {
		nameList = strings.Split(name, ",")
	}
	if len(nameList) > 0 {
		filter.Name = ""
	}
	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx))
	tx.Where(&filter)

	if len(nameList) > 0 {
		tx.Where("name in (?)", nameList)
	}
	if len(statuslist) > 0 {
		tx.Where("status in (?)", statuslist)
	}

	err := tx.Find(&results).Error
	if err != nil {
		ctx.Log.Error("error while fetching filtering cards", zap.Any("filter", filter), zap.Error(err))
		return nil, err
	}
	return results, nil
}

func (t *Card) GetExecCards(ctx *context.Context, req *models.Card, cardfilter *dtos.GECardFilter) ([]*models.Card, error) {
	var results []*models.Card
	statusList := []string{constants.CardStatusBreached, constants.CardStatusCreated, constants.CardStatusWarning}
	currentTime := time.Now().UTC()

	q := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Where("status in ?", statusList).
		Where("visible_time < ?", currentTime).
		Where(&cardfilter).Order("estimate ASC")

	if req.EscalatedTo != "" {
		q = q.Where("((escalated_to = ? OR escalated_by_id @> ?) OR assigned_to = ?) ", req.EscalatedTo, req.EscalatedById, req.AssignedTo)
	} else {
		q = q.Where("assigned_to = ?", req.AssignedTo)
	}
	err := q.Find(&results).Error
	if err != nil {
		ctx.Log.Error("Error while ExecCards in DB", zap.Any("filter", req), zap.Error(err))
		return nil, err
	}
	return results, nil
}

func (t *Card) GetAssignedTo(ctx *context.Context, filter *models.Card, statusList []string) (*models.Card, error) {
	var res *models.Card

	conn := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("assigned_to, escalated_to").Where(&filter)
	if len(statusList) > 0 {
		conn.Where("status in ?", statusList)
	}
	err := conn.Scan(&res).Error
	if err != nil {
		ctx.Log.Error("Error while Getting executive from DB", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (t *Card) GetShipmentsWithFilter(ctx *context.Context, isShipment bool, rep_ids []string, shipmentFilter *dtos.GEShipmentFilter, cardFilter *dtos.GECardFilter) (models.CountCards, error) {
	var results models.CountCards
	currentTime := time.Now().UTC()
	statusList := []string{constants.CardStatusBreached, constants.CardStatusCreated, constants.CardStatusWarning}
	query := ctx.DB.Debug().Table(t.getTable(ctx)).
		Select(fmt.Sprintf("%s.assigned_to,%s.assigned_to_name ,%s.status, COUNT(DISTINCT(%s.id)) as task_count", t.getTable(ctx), t.getTable(ctx), t.getTable(ctx), t.getTable(ctx)))

	if isShipment {
		query = query.Joins(fmt.Sprintf("JOIN %s ON %s.id::TEXT = %s.instance_id ::TEXT", t.getShipmentsTable(ctx), t.getShipmentsTable(ctx), t.getTable(ctx))).
			Where("instance_id != '' and instance_type='shipment'").
			Where(fmt.Sprintf("%s.consol_id IS NULL OR %s.consol_id = ?", t.getShipmentsTable(ctx), t.getShipmentsTable(ctx)), uuid.Nil)
	} else {
		query = query.Joins(fmt.Sprintf("JOIN %s ON %s.id::TEXT = %s.instance_id::TEXT ", t.getRfqsTable(ctx), t.getRfqsTable(ctx), t.getTable(ctx))).Where("instance_id != '' and instance_type='rfq'")
	}

	query = query.Where("assigned_to in ?", rep_ids).
		Where("visible_time < ?", currentTime).
		Where(fmt.Sprintf("%s.status in (?)", t.getTable(ctx)), statusList).
		Where(&cardFilter).
		Group(fmt.Sprintf("%s.assigned_to, %s.assigned_to_name, %s.status", t.getTable(ctx), t.getTable(ctx), t.getTable(ctx)))

	if shipmentFilter.Code != "" {
		query = query.Where(fmt.Sprintf("%s.code=?", t.getShipmentsTable(ctx)), shipmentFilter.Code)
	}

	if shipmentFilter.Pol != "" {
		query = query.Where(fmt.Sprintf("%s.pod=?", t.getShipmentsTable(ctx)), shipmentFilter.Pol)
	}

	if shipmentFilter.Pod != "" {
		query = query.Where(fmt.Sprintf("%s.pod=?", t.getShipmentsTable(ctx)), shipmentFilter.Pod)
	}

	if shipmentFilter.ShipmentNature != "" {
		query = query.Where(fmt.Sprintf("%s.shipment_nature=?", t.getShipmentsTable(ctx)), shipmentFilter.ShipmentNature)
	}

	if shipmentFilter.ShipmentType != "" {
		query = query.Where(fmt.Sprintf("%s.type = ?", t.getShipmentsTable(ctx)), shipmentFilter.ShipmentType)
	}

	err := query.Find(&results).Error
	if err != nil {
		ctx.Log.Error("Error while getting bulk  --DB", zap.Any("bookingFilter", shipmentFilter), zap.Error(err))
	}

	return results, nil
}

func (t *Card) GetCardswithNoBooks(ctx *context.Context, rep_ids []string, cardFilter *dtos.GECardFilter) (models.CountCards, error) {
	var results models.CountCards
	currentTime := time.Now().UTC()
	statusList := []string{constants.CardStatusBreached, constants.CardStatusCreated, constants.CardStatusWarning}
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

func (t *Card) GetLiveEscaltedCardsOrgtree(ctx *context.Context) (map[string]dtos.TaskCounts, error) {

	statusList := []string{constants.CardStatusBreached, constants.CardStatusCreated, constants.CardStatusWarning}
	rows, err := ctx.DB.Debug().Raw(`SELECT
    unnested_uuid AS escalated_by_id,
    COUNT(id) AS uuid_count
FROM
    (
        SELECT
            id,assigned_to,
            unnest(escalated_by_id) AS unnested_uuid
        FROM
            `+t.getTable(ctx)+`
        WHERE
            escalated = true
            AND status IN (?)
			AND instance_id NOT IN (
                SELECT id::TEXT
                FROM `+t.getShipmentsTable(ctx)+`
                WHERE consol_id is NOT NULL AND consol_id != ?
            )
    ) AS subquery
    where subquery.assigned_to =subquery.unnested_uuid
GROUP BY
    unnested_uuid
ORDER BY
    uuid_count DESC`, statusList, uuid.Nil).Rows()
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

func (t *Card) GetCardCounts(ctx *context.Context) (map[string]dtos.TaskCounts, error) {

	statusList := []string{constants.CardStatusBreached, constants.CardStatusCreated, constants.CardStatusWarning}
	rows, err := ctx.DB.Debug().Raw(`SELECT assigned_to,status
	FROM  `+t.getTable(ctx)+`
	WHERE  visible_time < Now()
	AND escalated = 'false' AND status IN (?)
	AND ( instance_id NOT IN (SELECT id::TEXT
						 FROM   `+t.getShipmentsTable(ctx)+`
						 WHERE consol_id is NOT NULL AND consol_id != ?)
						) `, statusList, uuid.Nil).Rows()
	if err != nil {
		ctx.Log.Error("Error while Getting team cards from DB", zap.Error(err))
	}
	defer rows.Close()
	taskCounts := make(map[string]dtos.TaskCounts)
	for rows.Next() {
		var assignedTo string
		var status string
		err = rows.Scan(&assignedTo, &status)
		if err != nil {
			return nil, err
		}

		counts := taskCounts[assignedTo]

		switch status {
		case constants.CardStatusCreated:
			counts.PendingTasks++
		case constants.CardStatusBreached:
			counts.BreachedTasks++
		case constants.CardStatusWarning:
			counts.ExpiringTasks++
		}
		taskCounts[assignedTo] = counts
	}
	return taskCounts, nil
}

func (c *Card) ReassignAllShipmentAssignedCard(ctx *context.Context, req *dtos.ReExecCard, instance_ids []string, aid string) ([]models.Card, error) {

	statusList := []string{globals.StatusBreached, dtos.ActionCreated, globals.StatusWarning}

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
	err := ctx.DB.Debug().Table(c.getTable(ctx)).Where("status in (?)", statusList).Where("instance_id in (?) and assigned_to = ?", instance_ids, aid).Updates(newValues).Error
	if err != nil {
		ctx.Log.Error("Error while reassigning card in DB", zap.Error(err))
	}

	var updatedCards []models.Card
	err = ctx.DB.Debug().Table(c.getTable(ctx)).
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

func (c *Card) ReassignCard(ctx *context.Context, req *dtos.ReExecCard) (*models.Card, error) {
	flowInstanceIds := req.CardIds
	statusList := []string{constants.CardStatusBreached, constants.CardStatusCreated, constants.CardStatusWarning}
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

	err := ctx.DB.Debug().
		Table(c.getTable(ctx)).
		Where("status IN (?)", statusList).
		Where("id IN (?)", flowInstanceIds).
		Updates(newValues).Error
	if err != nil {
		ctx.Log.Error("Error while reassigning card in DB", zap.Error(err))
		return nil, err
	}

	var updatedCard *models.Card

	err = ctx.DB.Debug().Table(c.getTable(ctx)).
		Where("status IN (?)", statusList).
		Where("id in (?)", flowInstanceIds).
		Find(&updatedCard).
		Error
	if err != nil {
		ctx.Log.Error("Error while retrieving updated cards from DB", zap.Error(err))
		return nil, err
	}

	return updatedCard, nil
}

func (c *Card) EscalateCard(ctx *context.Context, filter *dtos.ReExecCard, EscalatedByList pq.StringArray, id string) error {
	newValues := map[string]interface{}{
		"escalated":              true,
		"escalated_to":           filter.EscalateTo,
		"escalated_by_id":        EscalatedByList,
		"escalated_by":           filter.EscalatedByName,
		"escalation_reason":      filter.EscalationReason,
		"escalation_description": filter.EscalationRemarks,
		"updated_at":             time.Now().UTC(),
	}
	err := ctx.DB.Debug().Table(c.getTable(ctx)).Where("status not in ('Completed','Delete') and id = ?", id).Updates(newValues).Error
	if err != nil {
		ctx.Log.Error("Error while escalating card in DB", zap.Error(err))
	}

	return nil
}

func (c *Card) GetNonEscalatedCards(ctx *context.Context, filter *models.Card, statuslist []string, nonesc bool) ([]models.Card, error) {

	var results []models.Card
	query := ctx.DB.Debug().Table(c.getTable(ctx)).Where(&filter)
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

func (t *Card) GetDistinctColumnDetails(ctx *context.Context, label, value, cardType, querry string, ids, cardStatus []string) ([]dtos.CardLabelLists, error) {

	var results []dtos.CardLabelLists

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Distinct(" cards."+label+" AS label , "+"cards."+value+" AS value, '"+cardType+"' AS type ").Where("cards.assigned_to IN (?) and cards.status IN (?) and cards."+label+" ILIKE (?) AND cards.visible_time <= NOW()", ids, cardStatus, "%"+querry+"%").Find(&results).Error
	if err != nil {
		ctx.Log.Error("error while getting distinct column details", zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (t *Card) GetDistinctColumnDetailsInstanceData(ctx *context.Context, val dtos.CardLabelLists, filter *dtos.FilterCardLabel, ids, cardStatus []string) ([]dtos.CardLabelLists, error) {
	var results []dtos.CardLabelLists

	queryFunc := func(table, instanceType string) ([]dtos.CardLabelLists, error) {
		var queryResults []dtos.CardLabelLists
		query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table("cards").
			Select(fmt.Sprintf("DISTINCT %s.%s AS label, %s.%s AS value, ? AS type, ? AS for_type", table, val.Label, table, val.Value), val.Type, val.ForType).
			Joins(fmt.Sprintf("JOIN %s ON cards.instance_id IS NOT NULL AND cards.instance_type = ? AND cards.instance_id::text = %s.id::text", table, table), instanceType).
			Where(fmt.Sprintf("(cards.assigned_to IN (?) OR (escalated_to = ? OR escalated_by_id @> ?)) AND cards.status IN (?) AND %s.%s ILIKE ? AND cards.visible_time <= NOW()", table, val.Label), ids, filter.AdminID, pq.StringArray{filter.AdminID}, cardStatus, "%"+filter.Querry+"%").
			Find(&queryResults)

		if err := query.Error; err != nil {
			ctx.Log.Error(fmt.Sprintf("error while getting distinct column details for %s", table), zap.Error(err))
			return nil, err
		}
		return queryResults, nil
	}

	getTable := func(forType string) (string, string) {
		switch forType {
		case constants.WorkflowTypeRFQ:
			return t.getRfqsTable(ctx), constants.WorkflowTypeRFQ
		case constants.WorkflowTypeShipment:
			return t.getShipmentsTable(ctx), constants.WorkflowTypeShipment
		case constants.WorkflowTypeCONSOL:
			return t.getShipmentsTable(ctx), constants.WorkflowTypeCONSOL
		default:
			return "", ""
		}
	}

	if val.ForType == "All" {
		rfqResults, err := queryFunc(t.getRfqsTable(ctx), constants.WorkflowTypeRFQ)
		if err != nil {
			return nil, err
		}
		results = append(results, rfqResults...)

		shipmentResults, err := queryFunc(t.getShipmentsTable(ctx), constants.WorkflowTypeShipment)
		if err != nil {
			return nil, err
		}
		results = append(results, shipmentResults...)

		consolShipmentResults, err := queryFunc(t.getShipmentsTable(ctx), constants.WorkflowTypeCONSOL)
		if err != nil {
			return nil, err
		}
		results = append(results, consolShipmentResults...)

	} else {
		table, instanceType := getTable(val.ForType)
		if table == "" {
			return nil, fmt.Errorf("unknown ForType: %s", val.ForType)
		}

		queryResults, err := queryFunc(table, instanceType)
		if err != nil {
			return nil, err
		}
		results = append(results, queryResults...)
	}

	return results, nil
}

func (c *Card) GetCardsUsingMulti(ctx *context.Context, cardIds []string) ([]models.Card, error) {
	var results []models.Card
	statusList := []string{constants.CardStatusBreached, constants.CardStatusCreated, constants.CardStatusWarning}

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(c.getTable(ctx)).Where("status in (?)", statusList).Where("id in (?)", cardIds).Find(&results).Error
	if err != nil {
		ctx.Log.Error("Unable to get cards.", zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (c *Card) GetAllPendingCards(ctx *context.Context) ([]models.Card, error) {
	var results []models.Card
	excludedStatuses := []string{
		constants.CardStatusCompleted,
		constants.CardStatusDeleted,
		dtos.ActionDelete,
	}
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(c.getTable(ctx)).Where("status NOT IN (?)", excludedStatuses).Find(&results).Error
	if err != nil {
		ctx.Log.Error("error while getting pending cards.", zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (c *Card) UpdateStatus(ctx *context.Context, id, status string) error {

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(c.getTable(ctx)).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		ctx.Log.Error("error in updating card", zap.Any("card_id", id), zap.Any("status", status), zap.Error(err))
		return err
	}

	return nil
}

func (c *Card) GetCardsFiltered(ctx *context.Context, req *models.Card) ([]*models.Card, error) {
	var results []*models.Card

	q := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(c.getTable(ctx))

	if req.CompanyId != "" {
		q = q.Where("company_id = ?", req.CompanyId)
	}

	if req.Status != "" {

		statusList := strings.Split(req.Status, ",")
		q = q.Where("status in (?)", statusList)
	}

	if req.Name != "" {

		nameList := strings.Split(req.Name, ",")
		q = q.Where("name in (?)", nameList)
	}

	if req.FlowInstanceId != uuid.Nil {
		q = q.Where("flow_instance_id = ?", req.FlowInstanceId)
	}

	if req.InstanceId != "" {
		q = q.Where("instance_id = ?", req.InstanceId)
	}

	if req.InstanceType != "" {
		q = q.Where("instance_type = ?", req.InstanceType)
	}

	err := q.Find(&results).Error
	if err != nil {
		ctx.Log.Error("error while getting pending cards.", zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (c *Card) GetBulkCardsByCompanyId(ctx *context.Context, cids []string) ([]*models.Card, error) {

	var results []*models.Card

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(c.getTable(ctx)).Where("company_id IN (?)", cids).Find(&results).Error
	if err != nil {
		ctx.Log.Error("error while getting bulk cards.", zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (c *Card) UpsertCardById(ctx *context.Context, card *models.Card, action string) error {
	newValues := make(map[string]interface{})
	if action == constants.ActionDecline {
		newValues = map[string]interface{}{
			"escalated":       card.Escalated,
			"escalated_to":    card.EscalatedTo,
			"escalated_by_id": card.EscalatedById,
			"escalated_by":    card.EscalatedBy,
			"decline_reason":  card.DeclineReason,
		}
	}

	if action == constants.ActionExtendTimeline {
		newValues = map[string]interface{}{
			"estimate":        card.Estimate,
			"escalated":       card.Escalated,
			"escalated_to":    card.EscalatedTo,
			"escalated_by_id": card.EscalatedById,
			"escalated_by":    card.EscalatedBy,
			"status":          card.Status,
		}
	}
	card.UpdatedAt = time.Now().UTC()
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table("cards").Where("status not in ('Completed','Delete') and id = ?", card.Id).Updates(newValues).Error
	if err != nil {
		ctx.Log.Error("Error while updating card in DB", zap.Error(err))
		return err
	}
	return nil
}

func (c *Card) GetBulkCardsFiltered(ctx *context.Context, filters *models.BulkCardsFilters) ([]*models.Card, error) {
	var results []*models.Card

	q := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(c.getTable(ctx))

	if len(filters.InstanceIds) > 0 {
		q = q.Where("instance_id IN (?)", filters.InstanceIds)
	}

	if len(filters.Statuses) > 0 {
		q = q.Where("status IN (?)", filters.Statuses)
	}

	if len(filters.Names) > 0 {
		q = q.Where("name IN (?)", filters.Names)
	}

	err := q.Find(&results).Error
	if err != nil {
		ctx.Log.Error("error while getting pending cards.", zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (c *Card) DeleteBookingRequestByID(ctx *context.Context, id string, updatedAt time.Time) (bool, error) {

	tx := ctx.DB.Debug().WithContext(ctx.Request.Context()).
		Table(c.getTable(ctx)).
		Where("instance_id = ? AND status NOT IN ('Completed', 'Delete')", id).
		Updates(map[string]interface{}{
			"status":     constants.ActionDelete,
			"updated_at": updatedAt,
		})

	if tx.Error != nil {
		ctx.Log.Error("errorlog", zap.Error(tx.Error))
		return false, tx.Error
	}

	if tx.RowsAffected == 0 {
		ctx.Log.Error("no rows affected", zap.Any("Task Id", id))
		return false, nil
	}

	return true, nil
}

func (t *Card) BulkUpsert(ctx *context.Context, cardIds []string) error {

	if ctx == nil || ctx.DB == nil {
		return errors.New("invalid context or database connection")
	}
	if len(cardIds) == 0 {
		return nil
	}

	by := ""
	if ctx.Account != nil {
		by = ctx.Account.ID.String()
	}

	tx := ctx.DB.WithContext(ctx.Request.Context()).
		Table(t.getTable(ctx)).
		Where("id IN (?)", cardIds).
		Updates(map[string]interface{}{
			"status":       dtos.ActionDelete,
			"updated_at":   time.Now().UTC(),
			"updated_by":   by,
			"completed_at": time.Now().UTC(),
			"completed_by": by,
		})
	if tx.Error != nil {
		ctx.Log.Error("errorlog", zap.Error(tx.Error))
		return tx.Error
	}
	return nil

}
