package rfq

import (
	"strconv"
	"strings"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IRfq interface {
	Upsert(ctx *context.Context, m ...*models.Rfq) error
	UpdateRFQIsShipmentConverted(ctx *context.Context, rfqId uuid.UUID, value bool) error
	Get(ctx *context.Context, id string) (*models.Rfq, error)
	GetAll(ctx *context.Context, ids []string) ([]*models.Rfq, error)
	Delete(ctx *context.Context, id string) error
	CheckCode(ctx *context.Context, code string) (bool, error)
	GetRfqsCount(ctx *context.Context, m *dtos.GetRfqsFiltersReq) ([]*models.GetRfqCountRes, error)
	GetRfqsPaginated(ctx *context.Context, ms *dtos.GetRfqsFiltersReq) ([]*models.Rfq, error)
	GetSearchListForRfq(ctx *context.Context, m *dtos.GetRfqsFiltersReq) ([]*dtos.RfqFilter, error)
	Update(ctx *context.Context, m *models.Rfq) error
	UpdateWithTx(ctx *context.Context, tx *gorm.DB, m *models.Rfq) error
	GetCustomerDashboardEnquiryQuoteCount(ctx *context.Context, cids []string) (int64, int64, error)
	GetRfqCardFilter(ctx *context.Context, ids []string, req dtos.CardInstanceReq) ([]*models.Rfq, error)
	GetEnquiriesForPartner(ctx *context.Context, id string) ([]*models.Rfq, error)
	GetEnquiriesForCustomer(ctx *context.Context, cid string) ([]*models.Rfq, error)
	GetRfqByDealId(ctx *context.Context, dealId string) (*models.Rfq, error)
	GetLiveEnquiresCount(ctx *context.Context, filter *dtos.DashboardFilters) (int, error)
	GetLiveRFQsCount(ctx *context.Context, filter *dtos.DashboardFilters) (int, error)
	GetLiveEnquiries(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.LiveViewResponse, error)
	GetLiveQuotes(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.LiveViewResponse, error)
	GetFunnelViewCount(ctx *context.Context, filter *dtos.DashboardFilters) (int, error)
	GetforFunnelCountFilter(ctx *context.Context, filter *dtos.DashboardFilters) ([]string, error)
	GetforFunnelFilterQuotes(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.LiveViewResponse, error)
	GetQuoteActivities(ctx *context.Context, cid string) ([]*models.Rfq, error)
	GetPaginatedConsolRfqs(ctx *context.Context, req *dtos.ConsolGetReq) ([]*dtos.ConsolParameters, error)
	GetCountsPendingConsol(ctx *context.Context, req *dtos.ConsolGetReq) (int64, error)
}

type Rfq struct {
}

func NewRfq() IRfq {
	return &Rfq{}
}

func (t *Rfq) getTable(ctx *context.Context) string {
	if ctx.TenantID == "" {
		ctx.TenantID = "public"
	}
	return ctx.TenantID + "." + "rfqs"
}

func (t *Rfq) Upsert(ctx *context.Context, m ...*models.Rfq) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Save(m).Error
}

func (t *Rfq) UpdateWithTx(ctx *context.Context, tx *gorm.DB, m *models.Rfq) error {
	return tx.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Debug().Where("id = ?", m.ID).Updates(m).Error
}

func (t *Rfq) Update(ctx *context.Context, m *models.Rfq) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Updates(m).Error
}

func (t *Rfq) Get(ctx *context.Context, id string) (*models.Rfq, error) {
	var result models.Rfq
	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfq.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *Rfq) Delete(ctx *context.Context, id string) error {
	var result models.Rfq
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Delete(&result, "id = ?", id).Error
	if err != nil {
		ctx.Log.Error("Unable to delete rfq.", zap.Error(err))
		return err
	}

	return err
}

func (t *Rfq) GetAll(ctx *context.Context, ids []string) ([]*models.Rfq, error) {
	var result []*models.Rfq
	if len(ids) == 0 {
		err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Find(&result).Error
		if err != nil {
			ctx.Log.Error("Unable to get rfqs.", zap.Error(err))
			return nil, err
		}
		return result, err
	}
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfqs.", zap.Error(err))
		return nil, err
	}

	return result, err
}

func (t *Rfq) CheckCode(ctx *context.Context, code string) (bool, error) {
	var result *models.Rfq
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("code = ?", code).First(&result).Error
	return (err == nil && result.ID != uuid.Nil), err
}

func (t *Rfq) UpdateRFQIsShipmentConverted(ctx *context.Context, rfqId uuid.UUID, value bool) error {
	return ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id", rfqId).Update("is_shipement_converted", value).Error
}

func (t *Rfq) GetRfqsCount(ctx *context.Context, m *dtos.GetRfqsFiltersReq) ([]*models.GetRfqCountRes, error) {
	var result []*models.GetRfqCountRes

	q := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("rfqs.status", "count(DISTINCT rfqs.id) as count")
	reqStatus := m.Status
	m.Status = ""

	t.getRfqsFilterQuery(ctx, m, q)

	if reqStatus != "" {
		q = q.Where("rfqs.status = ? or rfqs.status = ?", reqStatus, constants.RfqStatusExpired)
	}

	err := q.Group("rfqs.status").Scan(&result).Error
	if err != nil {
		ctx.Log.Error("failed to get count of rfqs", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (t *Rfq) GetRfqsPaginated(ctx *context.Context, ms *dtos.GetRfqsFiltersReq) ([]*models.Rfq, error) {

	var res []*models.Rfq
	offset := int(ms.Count * (ms.Pg - 1))
	limit := int(ms.Count)

	q := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("DISTINCT rfqs.*")

	t.getRfqsFilterQuery(ctx, ms, q)

	if strings.EqualFold(ms.SortBy, constants.SortByCargoReadydate) {
		if strings.EqualFold(ms.SortOrder, constants.SortOrderAsc) {
			q = q.Order("rfqs.cargo_ready_date ASC")
		} else {
			q = q.Order("rfqs.cargo_ready_date DESC")
		}
	} else {
		q = q.Order("rfqs.created_at DESC")
	}

	q = q.Limit(limit).
		Offset(offset)

	err := q.Scan(&res).Error
	if err != nil {
		ctx.Log.Error("failed to get count of rfqs", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (t *Rfq) getRfqsFilterQuery(ctx *context.Context, m *dtos.GetRfqsFiltersReq, q *gorm.DB) {
	// Base conditions to exclude deleted RFQs and shipments already converted
	q.Where("(rfqs.is_deleted = false OR rfqs.is_deleted IS NULL) AND (rfqs.is_shipment_converted = false) AND rfqs.created_by != ? AND rfqs.company_id != ? AND rfqs.region_id != ?", "", uuid.Nil, uuid.Nil)

	// OR Queries for MyTab and Search
	orQueriesForMyTab := ctx.DB.Where("")
	orQueriesForMyTabPresent := false

	orQueriesForSearch := ctx.DB.Where("")
	orQueriesForSearchPresent := false

	// Prepare adminIds and customerIds with correct formatting using QuotedSlice
	adminIdsString := strings.Join(utils.QuotedSlice(m.AdminIds), ",")
	customerIdsString := strings.Join(utils.QuotedSlice(m.CustomerIds), ",")

	// cards table check (for MyTab)
	if len(m.AdminIds) > 0 {
		cardsQuery := ctx.DB.WithContext(ctx.Request.Context()).Table("cards").
			Select("DISTINCT instance_id").
			Where("assigned_to::text = ANY(ARRAY[" + adminIdsString + "])").
			Or("completed_by::text = ANY(ARRAY[" + adminIdsString + "])")

		orQueriesForMyTab.Or("rfqs.id::text IN (?)", cardsQuery)
		orQueriesForMyTabPresent = true
	}

	// Line Items Table Check (for MyTab)
	if len(m.AdminIds) > 0 {
		lineItemsQueryForMyTab := ctx.DB.WithContext(ctx.Request.Context()).Table("line_items").
			Select("DISTINCT quote_id").
			Where("updated_by::text = ANY(ARRAY[" + adminIdsString + "])")

		rfqQuotesQueryForMyTab := ctx.DB.WithContext(ctx.Request.Context()).Table("rfq_quotes").
			Select("DISTINCT rfq_id").
			Where("quote_id IN (?)", lineItemsQueryForMyTab)

		orQueriesForMyTab.Or("rfqs.id IN (?)", rfqQuotesQueryForMyTab)
		orQueriesForMyTabPresent = true
	}

	// Search: Handle partner_id
	if m.PartnerID != "" {
		// Apply Partner ID condition for line items
		lineItemsQueryForSearch := ctx.DB.WithContext(ctx.Request.Context()).Table("line_items").
			Select("DISTINCT quote_id").
			Where("partner_id::text = ?", m.PartnerID)

		rfqQuotesQueryForSearch := ctx.DB.WithContext(ctx.Request.Context()).Table("rfq_quotes").
			Select("DISTINCT rfq_id").
			Where("quote_id IN (?)", lineItemsQueryForSearch)

		// Apply the condition for the partner-related line items
		orQueriesForSearch.Or("rfqs.id IN (?)", rfqQuotesQueryForSearch)
		orQueriesForSearchPresent = true
	}

	// General RFQS Checks (for MyTab and Search)
	if len(m.AdminIds) > 0 {
		orQueriesForMyTab.Or("rfqs.sales_executive_id::text = ANY(ARRAY[" + adminIdsString + "])").
			Or("rfqs.created_by::text = ANY(ARRAY[" + adminIdsString + "])")
		orQueriesForMyTabPresent = true
	}

	if len(m.CustomerIds) > 0 {
		orQueriesForMyTab.Or("rfqs.company_id::text = ANY(ARRAY[" + customerIdsString + "])")
		orQueriesForMyTabPresent = true
	}

	// Apply OR queries if needed
	if orQueriesForMyTabPresent {
		q.Where(ctx.DB.Where(orQueriesForMyTab))
	}

	if orQueriesForSearchPresent {
		q.Where(ctx.DB.Where(orQueriesForSearch))
	}

	// Filter by Q if provided (search term for RFQ IDs)
	if m.Q != "" {
		q.Where("rfqs.id = ANY(?::uuid[])", pq.Array(strings.Split(m.Q, ",")))
	}

	// Account Region ID check for RFQs
	if ctx.Account.RegionID != "" {
		q.Where("? = ANY(ARRAY[rfqs.region_id, rfqs.origin_region_id, rfqs.dest_region_id])", ctx.Account.RegionID)
	}

	// Customer ID and Company ID filters
	if m.CompanyID != "" {
		q.Where("rfqs.company_id = ANY(?::uuid[])", pq.Array(strings.Split(m.CompanyID, ",")))
	}

	// Date Range Filters (created_at between From and To)
	if !m.From.IsZero() {
		q = q.Where("rfqs.created_at >= ?", m.From)
	}

	if !m.To.IsZero() {
		q = q.Where("rfqs.created_at <= ?", m.To)
	}

	// CreatedBy filter for booking or tenant-specific exec bookings
	if m.CreatedBy != "" && m.CreatedBy != constants.TypeAll {
		execBook := m.CreatedBy == ctx.TenantID
		q = q.Where("rfqs.is_exec_booked = ?", execBook)
	}

	// Category filter
	if m.Category != "" {
		q = q.Where("rfqs.type = ?", m.Category)
	}

	// Shipment nature filter
	if m.ShipmentNature != "" {
		q = q.Where("rfqs.shipment_nature = ?", m.ShipmentNature)
	}

	// Status filter: Active vs Expired
	if m.Status != "" {
		if m.Status == constants.Active {
			q = q.Where("rfqs.status != ?", constants.RfqStatusExpired)
		} else {
			q = q.Where("rfqs.status = ?", m.Status)
		}
	}

	// POL and POD filters
	if m.POL != "" {
		q = q.Where("rfqs.POL = ?", m.POL)
	}

	if m.POD != "" {
		q = q.Where("rfqs.POD = ?", m.POD)
	}

	// POL and POD Country filters
	if m.POLCountry != "" {
		q = q.Where("rfqs.pol_country = ?", m.POLCountry)
	}

	if m.PODCountry != "" {
		q = q.Where("rfqs.pod_country = ?", m.PODCountry)
	}

	// Entity filter
	if m.Entity != "" {
		q = q.Where("rfqs.entity = ?", m.Entity)
	}

	// Booking Country filter
	if m.BookingCountryID != "" {
		q = q.Where("rfqs.region_id = ?", m.BookingCountryID)
	}

	// Filter by specific RFQ IDs (if provided)
	if len(m.Ids) > 0 {
		q.Where("rfq_id IN (?)", m.Ids)
	}
}

func (t *Rfq) GetSearchListForRfq(ctx *context.Context, m *dtos.GetRfqsFiltersReq) ([]*dtos.RfqFilter, error) {
	var res []*dtos.RfqFilter
	err := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("id", "code").
		Where("(is_deleted = false OR is_deleted IS NULL) AND (is_shipment_converted = false) AND ? = ANY(ARRAY[region_id,origin_region_id,dest_region_id]) AND code ilike ?", ctx.Account.RegionID, m.Q+"%").Scan(&res).Error
	if err != nil {
		ctx.Log.Error("failed to get search list for rfq", zap.Error(err))
		return nil, err
	}
	return res, err
}

func (t *Rfq) GetCustomerDashboardEnquiryQuoteCount(ctx *context.Context, cids []string) (int64, int64, error) {
	var enquiryCount, quoteCount int64

	enquiryQuery := ctx.DB.WithContext(ctx.Request.Context()).
		Table(t.getTable(ctx)).
		Select("id").
		Where("is_deleted = false").
		Where("is_shipment_converted = false").
		Where("company_id IN (?)", cids).
		Not("status IN (?)", []string{constants.RfqStatusConfirmed, constants.RfqStatusExpired})

	err := enquiryQuery.Debug().Count(&enquiryCount).Error
	if err != nil {
		return 0, 0, err
	}

	quoteQuery := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("distinct(rfqs.id)").
		Joins("INNER JOIN rfq_quotes rq ON rfqs.id = rq.rfq_id").
		Joins("INNER JOIN quotes q ON rq.quote_id = q.id").
		Where("rfqs.company_id IN (?)", cids).
		Where("rfqs.is_shipment_converted  = false").
		Where("rfqs.status = ?", constants.RfqStatusConfirmed).
		Where("rq.rfq_id IS NOT NULL").
		Where("rfqs.is_deleted = false")

	quoteQuery = quoteQuery.Where("company_id IN (?)", cids)

	err = quoteQuery.Debug().Count(&quoteCount).Error
	if err != nil {
		return 0, 0, err
	}

	return enquiryCount, quoteCount, nil
}

func (t *Rfq) GetRfqCardFilter(ctx *context.Context, ids []string, req dtos.CardInstanceReq) ([]*models.Rfq, error) {
	var results []*models.Rfq
	query := ctx.DB.WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Where("id in (?) ", ids).Where(&req)
	err := query.Find(&results).Error
	if err != nil {
		ctx.Log.Error("Error while getting bulk  --DB", zap.Any("id", ids), zap.Error(err))
	}
	return results, nil
}

func (t *Rfq) GetEnquiriesForPartner(ctx *context.Context, id string) ([]*models.Rfq, error) {
	var rfqs []*models.Rfq
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Select("DISTINCT rfqs.id").
		Joins("INNER JOIN rfq_quotes AS rq ON rfqs.id = rq.rfq_id").
		Joins("INNER JOIN quotes AS q ON rq.quote_id = q.id").
		Joins("INNER JOIN line_items AS li ON li.quote_id = q.id").
		Where("li.partner_id = ?", id).
		Where("q.is_approved = ?", false).
		Where("rfqs.status != 'expired'").
		Where("rq.rfq_id IS NOT NULL")
	err := query.Find(&rfqs).Error
	if err != nil {
		return nil, err
	}
	return rfqs, nil
}

func (t *Rfq) GetEnquiriesForCustomer(ctx *context.Context, cid string) ([]*models.Rfq, error) {
	var rfqs []*models.Rfq
	enquiryQuery := ctx.DB.WithContext(ctx.Request.Context()).
		Table(t.getTable(ctx)).
		Where("is_deleted = false").
		Where("is_shipment_converted = false").
		Where("company_id = ?", cid).
		Not("status IN (?)", []string{constants.RfqStatusExpired})

	err := enquiryQuery.Debug().Find(&rfqs).Error
	if err != nil {
		return nil, err
	}
	return rfqs, nil
}

func (t *Rfq) GetRfqByDealId(ctx *context.Context, dealId string) (*models.Rfq, error) {

	var result models.Rfq

	err := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).First(&result, "deal_id = ?", dealId).Error
	if err != nil {
		ctx.Log.Error("Unable to get rfq.", zap.Error(err))
		return nil, err
	}

	return &result, err
}

func (t *Rfq) GetQuoteActivities(ctx *context.Context, cid string) ([]*models.Rfq, error) {
	var result []*models.Rfq
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)+" rfqs").
		Select("rfqs.id, rfqs.code, rfqs.pol, rfqs.pod, cards.updated_at").
		Joins("INNER JOIN cards ON rfqs.id::text = cards.instance_id").
		Where("rfqs.company_id = ?", cid).
		Where("rfqs.status = ?", globals.QuoteStatusConfirmed).
		Where("rfqs.is_deleted = ? and rfqs.is_shipment_converted = ?", false, false).
		Where("cards.status IN ?", []string{"Created", "Breached"}).
		Where("cards.name = ?", constants.MilestoneQuoteConfirmation).
		Where("cards.type = ?", constants.MilestoneShipmentCreated)

	err := query.Find(&result).Error
	if err != nil {
		ctx.Log.Error("Error in fetching rfq card details", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (t *Rfq) GetPaginatedConsolRfqs(ctx *context.Context, req *dtos.ConsolGetReq) ([]*dtos.ConsolParameters, error) {
	var results []*dtos.ConsolParameters
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)+" consol").
		Select(`
	consol.id AS consol_id, 
	rq.quote_id AS consol_quote_id,
	consol.code AS consol_code,
	rc.type AS container_info,
	consol.occupied_cbm AS occupied_volume,
	q.etd AS etd,
	consol.pol_name AS pol,
	consol.pod_name AS pod,
	consol.created_at AS created_on,
	'pending' AS status,
	rc.cbm AS total_volume`).
		Joins("JOIN rfq_quotes AS rq ON (consol.id = rq.rfq_id)").
		Joins("JOIN quotes AS q ON (rq.quote_id = q.id)").
		Joins("JOIN rfqs_containers AS rc ON consol.id = rc.rfq_id").
		Where("consol.region_id = ? and consol.type = ?  and consol.status = ?", req.RegionId, globals.BookingTypeCONSOL, globals.ConsolStatusBuyTBA)

	if req.Id != "" {
		query.Where("consol.id = ?", req.Id)
	}

	if req.Q != "" {
		query.Where("consol.code ilike (?)", req.Q+"%")
	}

	query.Order("consol.created_at DESC")

	size := config.Get().PageSize
	pg, err := strconv.Atoi(ctx.Query("pg"))
	if err != nil {
		ctx.Log.Error("unable to parse pageNo")
		return nil, err
	}

	offset := size * (pg - 1)
	if req.Q == "" {
		if offset > 0 {
			query.Offset(offset)
		}

		query.Limit(size)
	}

	err = query.Find(&results).Error
	if err != nil {
		ctx.Log.Error("error while getting consols", zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (t *Rfq) GetCountsPendingConsol(ctx *context.Context, req *dtos.ConsolGetReq) (int64, error) {
	q := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)+" rfq").
		Select("count(rfq.id)").
		Where("rfq.type = ? AND rfq.status = ? AND rfq.is_deleted = false", globals.BookingTypeCONSOL, globals.ConsolStatusBuyTBA).
		Where("rfq.region_id = ?", req.RegionId)

	if req.Id != "" {
		q = q.Where("rfq.id = ?", req.Id)
	}
	if req.Q != "" {
		q = q.Where("rfq.code = ?", req.Q)
	}
	if req.Pol != "" {
		q = q.Where("rfq.pol = ?", req.Pol)
	}

	var pendingCount int64
	err := q.Scan(&pendingCount).Error
	if err != nil {
		return 0, err
	}

	return pendingCount, nil
}
