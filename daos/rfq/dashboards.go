package rfq

import (
	"time"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
)

func (t *Rfq) GetforFunnelFilterQuotes(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.LiveViewResponse, error) {
	query := ctx.DB.Table(t.getTable(ctx)).Select(`rfqs.id ,rfqs.code, extract(epoch from rfqs.created_at)::INTEGER as created_at,
		rfqs.pol_name as origin_port, rfqs.pod_name as dest_port,  rfqs.company_name,
		(CASE WHEN rfqs.type != 'FCL' THEN round(rfqs.occupied_cbm,2) ELSE 0 END) as volume, 
		(CASE WHEN rfqs.type != 'FCL' THEN round(rfqs.occupied_weight,2) ELSE 0 END) as weight, 
		rfqs.type, rfqs.shipment_nature, rfqs.teus as teu,rfqs.company_id ::text AS customer_id
		`).
		Joins("LEFT JOIN cards AS bc ON (rfqs.id)::TEXT = bc.instance_id AND bc.name = 'Buy Rates' AND bc.assigned_to = ?", filter.RequestedBy).
		Where("is_deleted = false").
		Where("rfqs.is_shipment_converted = false")

	if filter.Dashboard == constants.Procurement {
		query = query.Where("bc.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("rfqs.sales_executive_id = ?", filter.RequestedBy)
	}

	if filter.From > 0 {
		query = query.Where("rfqs.created_at >= ?", time.Unix(filter.From, 0))
	}

	if filter.To > 0 {
		query = query.Where("rfqs.created_at <= ?", time.Unix(filter.To, 0))
	}

	if filter.Type != "" {
		if filter.Type == constants.ShipmentTypeFCL+constants.ShipmentTypeLCL {
			query = query.Where("(rfqs.type = ? OR rfqs.type = ?)", constants.ShipmentTypeFCL, constants.ShipmentTypeLCL)
		} else {
			query = query.Where("rfqs.type = ?", filter.Type)
		}
	}

	switch filter.ReportType {
	case constants.FilterEnquiry:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusBuyTBA)
	case constants.DropoffFilterBuyrateAdded:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusSellTBA)
	case constants.DropoffFilterSellrateAdded:
		query = query.Where("(rfqs.status not ilike ? AND rfqs.status not ilike ?)", "%"+constants.RfqStatusSellTBA, "%"+constants.RfqStatusBuyTBA)
	}

	query = query.Order("created_at asc")
	if filter.Pg != -1 {
		var limit, offset int
		offset = (config.Get().PageSize) * (filter.Pg - 1)
		limit = (config.Get().PageSize)
		if offset > 0 {
			query = query.Offset(offset)
		}

		query = query.Limit(limit)
	}

	res := &dtos.LiveViewResponse{
		Data: []*dtos.LiveViewData{},
	}

	err := query.Debug().Find(&res.Data).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *Rfq) GetforFunnelCountFilter(ctx *context.Context, filter *dtos.DashboardFilters) ([]string, error) {
	var res []string
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Joins("LEFT JOIN cards AS bc ON (rfqs.id)::TEXT = bc.instance_id AND bc.name = 'Buy Rates' AND bc.assigned_to = ?", filter.RequestedBy).Where("is_deleted = false").
		Select("distinct(rfqs.type) as type").
		Where("rfqs.is_shipment_converted = false")

	if filter.Dashboard == constants.Procurement {
		query = query.Where("bc.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("rfqs.sales_executive_id = ?", filter.RequestedBy)
	}

	if filter.From > 0 {
		query = query.Where("rfqs.created_at >= ?", time.Unix(filter.From, 0))
	}

	if filter.To > 0 {
		query = query.Where("rfqs.created_at <= ?", time.Unix(filter.To, 0))
	}

	err := query.Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *Rfq) GetFunnelViewCount(ctx *context.Context, filter *dtos.DashboardFilters) (int, error) {
	var cnt int
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Joins("LEFT JOIN cards AS bc ON (rfqs.id)::TEXT = bc.instance_id AND bc.name = 'Buy Rates' AND bc.assigned_to = ?", filter.RequestedBy).Where("is_deleted = false").
		Where("rfqs.is_shipment_converted = false")

	switch filter.MetricsFilter {
	case "volume":
		query = query.Select("COALESCE(sum(rfqs.occupied_cbm), 0)::INT as volume")
	case "weight":
		query = query.Select("COALESCE(sum(rfqs.occupied_weight), 0)::INT as weight")
	case "teus":
		query = query.Select("count(rfqs.id)")
	default:
		query = query.Select("count(rfqs.id)")
	}

	if filter.Dashboard == constants.Procurement {
		query = query.Where("bc.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("rfqs.sales_executive_id = ?", filter.RequestedBy)
	}

	if filter.Type == constants.ShipmentTypeFCL || filter.Type == constants.ShipmentTypeLCL || filter.Type == constants.ShipmentTypeAIR {
		if filter.Type == constants.ShipmentTypeFCL+constants.ShipmentTypeLCL {
			query = query.Where("(rfqs.type = ? OR rfqs.type = ?)", constants.ShipmentTypeFCL, constants.ShipmentTypeLCL)
		} else {
			query = query.Where("rfqs.type = ?", filter.Type)
		}
	}

	if filter.From > 0 {
		query = query.Where("rfqs.created_at >= ?", time.Unix(filter.From, 0))
	}

	if filter.To > 0 {
		query = query.Where("rfqs.created_at <= ?", time.Unix(filter.To, 0))
	}

	switch filter.ReportType {
	case constants.FilterEnquiry:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusBuyTBA)
	case constants.FilterBuyrateAdded:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusSellTBA)
	case constants.FilterSellrateAdded:
		query = query.Where("(rfqs.status not ilike ? AND rfqs.status not ilike ?)", "%"+constants.RfqStatusSellTBA, "%"+constants.RfqStatusBuyTBA)
	}

	err := query.Scan(&cnt).Error
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

func (t *Rfq) GetLiveQuotes(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.LiveViewResponse, error) {
	query := ctx.DB.Table(t.getTable(ctx)).Select(`rfqs.id ,rfqs.code, extract(epoch from rfqs.created_at)::INTEGER as created_at,
		rfqs.pol_name as origin_port, rfqs.pod_name as dest_port,  rfqs.company_name,
		(CASE WHEN rfqs.type != 'FCL' THEN round(rfqs.occupied_cbm,2) ELSE 0 END) as volume, 
		(CASE WHEN rfqs.type != 'FCL' THEN round(rfqs.occupied_weight,2) ELSE 0 END) as weight,
		 rfqs.type, rfqs.shipment_nature, rfqs.teus as teu,rfqs.company_id ::text AS customer_id
		`).
		Joins("JOIN rfq_quotes AS rq ON (rfqs.id = rq.rfq_id)").
		Joins("JOIN quotes AS q ON (rq.quote_id = q.id)").
		Joins("LEFT JOIN cards AS bc ON (rfqs.id)::TEXT = bc.instance_id AND bc.name = 'Buy Rates' AND bc.assigned_to = ?", filter.RequestedBy).
		Where("is_deleted = false").
		Where("q.is_approved = true").
		Where("rfqs.is_shipment_converted = false")

	if filter.Dashboard == constants.Procurement {
		query = query.Where("bc.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("rfqs.sales_executive_id = ?", filter.RequestedBy)
	}

	switch filter.ReportType {
	case constants.FilterLostQuotes:
		query = query.Where("rfqs.status = ?", constants.RfqStatusExpired).Where("rfqs.expires_at > ?", time.Now().Add(-1*24*time.Hour*time.Duration(constants.LiveDisplayDays)))
	case constants.FilterNearingExpiry:
		query = query.Where("rfqs.status != ?", constants.RfqStatusExpired).Where("rfqs.expires_at < ?", time.Now().Add(1*time.Hour*time.Duration(config.Get().LivequotesDays)))
	case constants.FilterBookingPending:
		query = query.Where("rfqs.status != ?", constants.RfqStatusExpired).Where("rfqs.expires_at >= ?", time.Now().Add(1*time.Hour*time.Duration(config.Get().LivequotesDays)))
	}

	query = query.Order("created_at asc")
	if filter.Pg != -1 {
		var limit, offset int
		offset = (config.Get().PageSize) * (filter.Pg - 1)
		limit = (config.Get().PageSize)
		if offset > 0 {
			query = query.Offset(offset)
		}

		query = query.Limit(limit)
	}

	res := &dtos.LiveViewResponse{
		Data: []*dtos.LiveViewData{},
	}

	err := query.Debug().Find(&res.Data).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *Rfq) GetLiveEnquiries(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.LiveViewResponse, error) {
	query := ctx.DB.Table(t.getTable(ctx)).
		Joins("LEFT JOIN cards AS bc ON (rfqs.id)::TEXT = bc.instance_id AND bc.name = 'Buy Rates' AND bc.assigned_to = ? AND bc.status != 'Delete'", filter.RequestedBy).
		Joins("LEFT JOIN cards AS sc ON (rfqs.id)::TEXT = sc.instance_id AND sc.name = 'Sell Rates' AND sc.assigned_to = ? AND sc.status != 'Delete'", filter.RequestedBy).
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Buy Rates' AND status != 'Delete' GROUP BY instance_id) as bre ON ((rfqs.id)::TEXT = bre.instance_id)").
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Sell Rates' AND status != 'Delete' GROUP BY instance_id) as sre ON ((rfqs.id)::TEXT = sre.instance_id)").
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Quote Approval' AND status != 'Delete' GROUP BY instance_id) as ae ON ((rfqs.id)::TEXT = ae.instance_id)").
		Where("is_deleted = false").
		Where("rfqs.is_shipment_converted = false")

	if filter.Dashboard == constants.Procurement {
		query = query.Where("bc.name = ?", constants.CardBuyRate)

		query = query.Select(`rfqs.id ,rfqs.code, extract(epoch from rfqs.created_at)::INTEGER as created_at,
		rfqs.pol_name as origin_port, rfqs.pod_name as dest_port,  rfqs.company_name,
		(CASE WHEN rfqs.type != 'FCL' THEN round(rfqs.occupied_cbm,2) ELSE 0 END) as volume, 
		(CASE WHEN rfqs.type != 'FCL' THEN round(rfqs.occupied_weight,2) ELSE 0 END) as weight, 
		rfqs.type, rfqs.shipment_nature, rfqs.teus as teu,
		CASE WHEN (bc.completed_at > bc.estimate OR (now() > bc.estimate and bc.status != 'Completed')) THEN true ELSE false END as is_breached,
		rfqs.company_id ::text AS customer_id`)
	}

	if filter.Dashboard == constants.Procurement && filter.ReportType == constants.CardBuyRate {
		query = query.Where("bc.status IS DISTINCT FROM 'Completed'")
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("rfqs.sales_executive_id = ?", filter.RequestedBy)

		query = query.Select(`rfqs.id ,rfqs.code, extract(epoch from rfqs.created_at)::INTEGER as created_at,
		rfqs.pol_name as origin_port, rfqs.pod_name as dest_port,  rfqs.company_name,
		(CASE WHEN rfqs.type != 'FCL' THEN round(rfqs.occupied_cbm,2) ELSE 0 END) as volume, 
		(CASE WHEN rfqs.type != 'FCL' THEN round(rfqs.occupied_weight,2) ELSE 0 END) as weight,
		rfqs.type, rfqs.shipment_nature, rfqs.teus as teu,
		CASE WHEN (sc.completed_at > sc.estimate OR (now() > sc.estimate and sc.status != 'Completed')) THEN true ELSE false END as is_breached,
		rfqs.company_id ::text AS customer_id`)
	}

	switch filter.ReportType {
	case constants.FilterEnquiry:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusBuyTBA)
	case constants.DropoffFilterBuyrateAdded:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusSellTBA)
	case constants.DropoffFilterSellrateAdded:
		query = query.Where("(rfqs.status not ilike ? AND rfqs.status not ilike ?)", "%"+constants.RfqStatusSellTBA, "%"+constants.RfqStatusBuyTBA)
	}

	switch filter.ReportType {
	case constants.LiveLostEnquiries:
		query = query.Where("rfqs.status = ?", constants.RfqStatusExpired).Where("rfqs.expires_at > ?", time.Now().Add(-1*24*time.Hour*time.Duration(constants.LiveDisplayDays)))
	case constants.LiveBuyRatePending:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusBuyTBA).Where("(bre.expired = false OR bre.expired IS NULL)")
	case constants.LiveBuyRateExpired:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusBuyTBA).Where("bre.expired = true")
	case constants.LiveSellRatePending:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusSellTBA).Where("(sre.expired = false OR sre.expired IS NULL)")
	case constants.LiveSellRateExpired:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusSellTBA).Where("sre.expired = true")
	case constants.LiveApprovalPending:
		query = query.Where("rfqs.status = ?", globals.QuoteStatusPricingApprovalPending).Where("(ae.expired = false OR ae.expired IS NULL)")
	case constants.LiveApprovalExpired:
		query = query.Where("rfqs.status = ?", globals.QuoteStatusPricingApprovalPending).Where("ae.expired = true")
	case constants.LiveBuyRate:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusBuyTBA)
	case constants.LiveSellRate:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusSellTBA)
	case constants.CardApproval:
		query = query.Where("rfqs.status = ?", globals.QuoteStatusPricingApprovalPending)
	}

	query = query.Order("created_at asc")
	if filter.Pg != -1 {
		var limit, offset int
		offset = (config.Get().PageSize) * (filter.Pg - 1)
		limit = (config.Get().PageSize)
		if offset > 0 {
			query = query.Offset(offset)
		}

		query = query.Limit(limit)
	}

	res := &dtos.LiveViewResponse{
		Data: []*dtos.LiveViewData{},
	}

	err := query.Debug().Find(&res.Data).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *Rfq) GetLiveEnquiresCount(ctx *context.Context, filter *dtos.DashboardFilters) (int, error) {
	var cnt int
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Select("count(rfqs.id)").Table(t.getTable(ctx)).
		Joins("LEFT JOIN cards AS bc ON (rfqs.id)::TEXT = bc.instance_id AND bc.name = 'Buy Rates' AND bc.assigned_to = ?", filter.RequestedBy).
		Joins("LEFT JOIN cards AS sc ON (rfqs.id)::TEXT = sc.instance_id AND sc.name = 'Sell Rates' AND sc.assigned_to = ?", filter.RequestedBy).
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Buy Rates' GROUP BY instance_id) as bre ON ((rfqs.id)::TEXT = bre.instance_id)").
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Sell Rates'  GROUP BY instance_id) as sre ON ((rfqs.id)::TEXT = sre.instance_id)").
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Quote Approval'  GROUP BY instance_id) as ae ON ((rfqs.id)::TEXT = ae.instance_id)").
		Where("is_deleted = false").
		Where("rfqs.is_shipment_converted = false")

	if filter.Dashboard == constants.Procurement {
		query = query.Where("bc.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("rfqs.sales_executive_id = ?", filter.RequestedBy)
	}

	switch filter.ReportType {
	case constants.QuickActionLostEnquiries:
		query = query.Where("rfqs.status = ?", constants.RfqStatusExpired).Where("rfqs.expires_at > ?", time.Now().Add(-1*24*time.Hour*time.Duration(constants.LiveDisplayDays)))
	case constants.QuickActionBuyPending:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusBuyTBA).Where("(bre.expired = false OR bre.expired IS NULL)")
	case constants.QuickActionBuyTATExpired:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusBuyTBA).Where("bre.expired = true")
	case constants.QuickActionSellPending:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusSellTBA).Where("(sre.expired = false OR sre.expired IS NULL)")
	case constants.QuickActionSellTATExpired:
		query = query.Where("rfqs.status ilike ?", "%"+constants.RfqStatusSellTBA).Where("sre.expired = true")
	case constants.QuickActionApprovalPending:
		query = query.Where("rfqs.status = ?", globals.QuoteStatusPricingApprovalPending).Where("(ae.expired = false OR ae.expired IS NULL)")
	case constants.QuickActionApprovalTATExpired:
		query = query.Where("rfqs.status = ?", globals.QuoteStatusPricingApprovalPending).Where("ae.expired = true")
	}

	err := query.Scan(&cnt).Error
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

func (t *Rfq) GetLiveRFQsCount(ctx *context.Context, filter *dtos.DashboardFilters) (int, error) {
	var cnt int
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Select("count(rfqs.id)").Table(t.getTable(ctx)).
		Joins("JOIN rfq_quotes AS rq ON (rfqs.id = rq.rfq_id)").
		Joins("JOIN quotes AS q ON (rq.quote_id = q.id)").
		Joins("LEFT JOIN cards AS bc ON (rfqs.id)::TEXT = bc.instance_id AND bc.name = 'Buy Rates' AND bc.assigned_to = ?", filter.RequestedBy).Where("is_deleted = false").Where("q.is_approved = true").
		Where("rfqs.is_shipment_converted = false")

	if filter.Dashboard == constants.Procurement {
		query = query.Where("bc.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("rfqs.sales_executive_id = ?", filter.RequestedBy)
	}

	switch filter.ReportType {
	case constants.FilterLostQuotes:
		query = query.Where("rfqs.status = ?", globals.QuoteStatusExpired).Where("rfqs.expires_at > ?", time.Now().Add(-1*24*time.Hour*time.Duration(constants.LiveDisplayDays)))
	case constants.FilterNearingExpiry:
		query = query.Where("rfqs.status != ?", constants.RfqStatusExpired).Where("rfqs.expires_at < ?", time.Now().Add(1*time.Hour*time.Duration(config.Get().LivequotesDays)))
	case constants.FilterBookingPending:
		query = query.Where("rfqs.status != ?", constants.RfqStatusExpired).Where("rfqs.expires_at >= ?", time.Now().Add(1*time.Hour*time.Duration(config.Get().LivequotesDays)))
	}

	err := query.Scan(&cnt).Error
	if err != nil {
		return 0, err
	}

	return cnt, nil
}
