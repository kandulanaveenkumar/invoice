package shipment

import (
	"strconv"
	"strings"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (t *Shipment) getListingFiltersQuery(ctx *context.Context, q *gorm.DB, isForCounts bool, adminIds, customerIds []string) {

	// OR Queries for MyTab and Search
	orQueriesForMyTab := ctx.DB.Where("")
	orQueriesForMyTabPresent := false

	orQueriesForSearch := ctx.DB.Where("")
	orQueriesForSearchPresent := false

	// Prepare adminIds and customerIds with correct formatting using QuotedSlice
	adminIdsString := strings.Join(utils.QuotedSlice(adminIds), ",")
	customerIdsString := strings.Join(utils.QuotedSlice(customerIds), ",")

	// cards table check (for MyTab)
	if len(adminIds) > 0 {
		cardsQuery := ctx.DB.WithContext(ctx.Request.Context()).Table("cards").
			Select("DISTINCT instance_id").
			Where("assigned_to::text = ANY(ARRAY[" + adminIdsString + "])").
			Or("completed_by::text = ANY(ARRAY[" + adminIdsString + "])")

		orQueriesForMyTab.Or("shipments.id::text IN (?)", cardsQuery)
		orQueriesForMyTabPresent = true
	}

	// Line Items Table Check (for MyTab)
	if len(adminIds) > 0 {
		lineItemsQueryForMyTab := ctx.DB.WithContext(ctx.Request.Context()).Table("line_items").
			Select("DISTINCT quote_id").
			Where("updated_by::text = ANY(ARRAY[" + adminIdsString + "])")

		orQueriesForMyTab.Or("shipments.quote_id IN (?)", lineItemsQueryForMyTab)
		orQueriesForMyTabPresent = true
	}

	// Search: Handle partner_id and partner search using "q" and "q_type"
	if ctx.Query("partner_id") != "" || (ctx.Query("q") != "" && ctx.Query("q_type") == "partner") {
		lineItemsQueryForSearch := ctx.DB.WithContext(ctx.Request.Context()).Table("line_items").
			Select("DISTINCT quote_id")

		// If partner_id is provided, apply the filter
		if ctx.Query("partner_id") != "" {
			lineItemsQueryForSearch.Where("partner_id::text = ANY(?)", pq.Array(strings.Split(ctx.Query("partner_id"), ",")))
		}

		// If q and q_type are provided for partner search, apply the filter for an array match
		if ctx.Query("q") != "" && ctx.Query("q_type") == "partner" {
			lineItemsQueryForSearch.Where("partner_id::text = ANY(?)", pq.Array(strings.Split(ctx.Query("q"), ",")))
		}

		// Apply the condition for the partner-related line items
		orQueriesForSearch.Or("shipments.quote_id IN (?)", lineItemsQueryForSearch)
		orQueriesForSearchPresent = true
	}

	// Shipment Parties Table Checks (for Search)
	if ctx.Query("q") != "" && ctx.Query("q_type") == "contact" {
		shipmentPartiesQuery := ctx.DB.WithContext(ctx.Request.Context()).Table("shipment_parties").
			Select("DISTINCT shipment_id").
			Where("address_id::text = ANY(?) AND region_id = ? AND type IN ('shipper', 'consignee')", pq.Array(strings.Split(ctx.Query("q"), ",")), ctx.Account.RegionID)

		orQueriesForSearch.Or("shipments.id IN (?)", shipmentPartiesQuery)
		orQueriesForSearchPresent = true
	}

	// General Shipment Checks (for MyTab and Search)
	if len(adminIds) > 0 {
		orQueriesForMyTab.Or("shipments.sales_executive_id::text = ANY(ARRAY[" + adminIdsString + "])").
			Or("shipments.created_by::text = ANY(ARRAY[" + adminIdsString + "])")
		orQueriesForMyTabPresent = true
	}

	if len(customerIds) > 0 {
		orQueriesForMyTab.Or("shipments.company_id::text = ANY(ARRAY[" + customerIdsString + "])")
		orQueriesForMyTabPresent = true
	}

	// Apply OR queries if needed
	if orQueriesForMyTabPresent {
		q.Where(ctx.DB.Where(orQueriesForMyTab))
	}

	if orQueriesForSearchPresent {
		q.Where(ctx.DB.Where(orQueriesForSearch))
	}

	// AND Queries (for additional filtering)
	q.Where("shipments.is_deleted != ?", true)

	if ctx.Query("q") != "" && ctx.Query("q_type") == "shipment" {
		q.Where("shipments.id::text = ANY(?)", pq.Array(strings.Split(ctx.Query("q"), ",")))
	}

	if ctx.Query("q") != "" && ctx.Query("q_type") == "customer" {
		q.Where("shipments.company_id::text = ANY(?)", pq.Array(strings.Split(ctx.Query("q"), ",")))
	}

	if ctx.Query("type") != "" {
		q.Where("shipments.type = ?", ctx.Query("type"))
	}

	if ctx.Query("from") != "" {
		q.Where("shipments.created_at >= ?", ctx.Query("from"))
	}

	if ctx.Query("to") != "" {
		q.Where("shipments.created_at <= ?", ctx.Query("to"))
	}

	if ctx.Query("company_id") != "" {
		q.Where("shipments.company_id = ?", ctx.Query("company_id"))
	}

	if ctx.Query("shipment_nature") != "" {
		q.Where("shipment_nature = ?", ctx.Query("shipment_nature"))
	}

	if ctx.Query("pol") != "" {
		q.Where("pol = ?", ctx.Query("pol"))
	}

	if ctx.Query("pod") != "" {
		q.Where("pod = ?", ctx.Query("pod"))
	}

	if ctx.Query("pol_country") != "" {
		q.Where("pol_country = ?", ctx.Query("pol_country"))
	}

	if ctx.Query("pod_country") != "" {
		q.Where("pod_country = ?", ctx.Query("pod_country"))
	}

	if ctx.Query("requested_by") != "" {
		q.Where("shipments.created_by = ?", ctx.Query("requested_by"))
	}

	if ctx.Query("consol_id") != "" {
		q.Where("consol_id = ?", ctx.Query("consol_id"))
	}

	if ctx.Query("entity") != "" {
		q.Where("entity = ?", ctx.Query("entity"))
	}

	if ctx.Query("booking_country_id") != "" {
		q.Where("shipments.region_id = ?", ctx.Query("booking_country_id"))
	}

	if ctx.Query("icd_pol") != "" {
		q.Where("origin_por = ?", ctx.Query("pol"))
	}

	if ctx.Query("icd_pod") != "" {
		q.Where("dest_por = ?", ctx.Query("pod"))
	}

	if ctx.Query("nominated_by") != "" {
		switch ctx.Query("nominated_by") {
		case constants.FreeHand:
			q.Where("is_agent_nominated = ? AND is_nomination = ?", false, false)
		case constants.AgentNomination:
			q.Where("is_agent_nominated = ?", true)
		case constants.WizNomination:
			q.Where("is_nomination = ?", true)
		}
	}

	if ctx.Account.RegionID != "" {
		q.Where("? = ANY(ARRAY[shipments.region_id, origin_region_id, dest_region_id])", ctx.Account.RegionID)
	}

	if ctx.Query("booking_id") != "" {
		q.Where("shipments.id IN (?)", strings.Split(ctx.Query("ids"), ","))
	}

	// Status Filters
	statuses := []string{}
	if ctx.Query("status") != "" {
		statuses = strings.Split(ctx.Query("status"), ",")
	}

	if !isForCounts && len(statuses) > 0 {
		statusFilter := ctx.DB.Where("")
		for i := range statuses {
			condition := "="
			if strings.HasPrefix(statuses[i], "!") {
				condition = "!="
			}
			if i == 0 {
				statusFilter = statusFilter.Where("shipments.status "+condition+" ?", strings.ReplaceAll(statuses[i], "!", ""))
				continue
			}
			statusFilter = statusFilter.Or("shipments.status "+condition+" ?", strings.ReplaceAll(statuses[i], "!", ""))
		}
		q.Where(statusFilter)
	}
}

func (t *Shipment) GetShipmentCounts(ctx *context.Context, adminIds, customerIds []string) ([]*models.ShipmentCount, error) {
	var counts []*models.ShipmentCount

	if ctx.Query("partner_id") != "" || ctx.Query("q_type") == "partner" {
		counts, err := t.GetShipmentCountsForPartner(ctx, adminIds)
		if err != nil {
			ctx.Log.Error("error while fetching GetShipmentCounts for partner", zap.Error(err))
			return nil, err
		}
		return counts, nil
	}

	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("shipments.status", "count(shipments.id) as count")
	t.getListingFiltersQuery(ctx, query, true, adminIds, customerIds)
	query.Group("shipments.status")

	err := query.Find(&counts).Error
	if err != nil {
		ctx.Log.Error("error while fetching GetShipmentCounts", zap.Error(err))
		return nil, err
	}

	return counts, nil
}

func (t *Shipment) GetShipmentList(ctx *context.Context, adminIds, customerIds []string) ([]*models.ShipmentList, error) {
	var shipments []*models.ShipmentList

	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("shipments.*", "quotes.eta", "quotes.etd").
		Joins("LEFT JOIN " + ctx.TenantID + ".quotes ON shipments.quote_id = quotes.id")
	t.getListingFiltersQuery(ctx, query, false, adminIds, customerIds)

	sortOrder := ` DESC `
	if ctx.Query("sort_order") == "asc" {
		sortOrder = ` ASC `
	}
	switch ctx.Query("sort_by") {
	case "eta":
		query.Order("eta " + sortOrder)
	case "etd":
		query.Order("etd " + sortOrder)
	case "cargo_ready":
		query.Order("cargo_ready_date " + sortOrder)
	default:
		query.Order("shipments.created_at " + sortOrder)
	}
	if ctx.Query("sort_by") != constants.BookedOn && ctx.Query("sort_by") != "" {
		query.Order("shipments.created_at" + sortOrder)
	}

	size, err := strconv.Atoi(ctx.Query("count"))
	if err != nil {
		ctx.Log.Error("unable to parse pageSize")
		return nil, err
	}
	pg, err := strconv.Atoi(ctx.Query("pg"))
	if err != nil {
		ctx.Log.Error("unable to parse pageNo")
		return nil, err
	}

	offset := size * (pg - 1)
	if offset > 0 {
		query.Offset(offset)
	}
	query.Limit(size)

	err = query.Find(&shipments).Error
	if err != nil {
		ctx.Log.Error("error while fetching GetShipmentCounts", zap.Error(err))
		return nil, err
	}

	return shipments, nil
}

func (t *Shipment) GetShipmentListSearch(ctx *context.Context) ([]*models.ShipmentSearchFilter, error) {
	var shipmentSearch []*models.ShipmentSearchFilter

	subQuery := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("id", "unnest(house_bill_nos) AS name", "created_at", "'shipment' AS type").Where("is_deleted = ? AND ? = ANY(ARRAY[region_id,origin_region_id,dest_region_id])", false, ctx.Account.RegionID)
	query1 := ctx.DB.WithContext(ctx.Request.Context()).Raw("WITH house_bills AS ( ? ) SELECT id, name, created_at, 'shipment' AS type FROM house_bills WHERE name ilike ?", subQuery, ctx.Query("q")+"%")

	query2 := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("id", "master_bill_no AS name", "created_at", "'shipment' AS type").
		Where("is_deleted = ? AND ? = ANY(ARRAY[region_id,origin_region_id,dest_region_id]) AND master_bill_no ilike ?", false, ctx.Account.RegionID, ctx.Query("q")+"%")

	query3 := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("id", "code AS name", "created_at", "'shipment' AS type").
		Where("is_deleted = ? AND ? = ANY(ARRAY[region_id,origin_region_id,dest_region_id]) AND code ilike ?", false, ctx.Account.RegionID, ctx.Query("q")+"%")

	err := ctx.DB.WithContext(ctx.Request.Context()).Raw("? UNION ? UNION ?", query1, query2, query3).Find(&shipmentSearch).Error
	if err != nil {
		ctx.Log.Error("error while fetching GetShipmentCounts", zap.Error(err))
		return nil, err
	}

	return shipmentSearch, nil
}

func (t *Shipment) GetShipmentCountsForPartner(ctx *context.Context, adminIds []string) ([]*models.ShipmentCount, error) {
	var partnerShipments []*models.ShipmentList
	var counts []*models.ShipmentCount
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("shipments.*", "quotes.eta", "quotes.etd").
		Joins("LEFT JOIN " + ctx.TenantID + ".quotes ON shipments.quote_id = quotes.id")
	t.getListingFiltersQuery(ctx, query, true, adminIds, nil)

	query.Group("shipments.id")
	query.Group("quotes.eta")
	query.Group("quotes.etd")
	err := query.Find(&partnerShipments).Error
	if err != nil {
		ctx.Log.Error("error while fetching GetShipmentCounts", zap.Error(err))
		return nil, err
	}
	statusMap := make(map[string]int64)
	for _, partnerShipment := range partnerShipments {
		_, statusExists := statusMap[partnerShipment.Status]
		if statusExists {
			statusMap[partnerShipment.Status] += 1
		} else {
			statusMap[partnerShipment.Status] = 1
		}
	}

	for status, statusCount := range statusMap {
		counts = append(counts, &models.ShipmentCount{
			Status: status,
			Count:  statusCount,
		})
	}
	return counts, nil

}
