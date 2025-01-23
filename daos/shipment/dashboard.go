package shipment

import (
	"time"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
)

func (t *Shipment) GetInsightDetailCompanies(ctx *context.Context, filter *dtos.DashboardFilters, ports []string) ([]*dtos.Company, error) {
	res := []*dtos.Company{}
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Select(`distinct shipments.company_id as id, shipments.company_name as name`).
		Joins("JOIN rfqs as r ON (r.id = shipments.rfq_id)").
		Joins("LEFT JOIN cards as c ON ((r.id)::TEXT = c.instance_id  AND c.assigned_to = ?)", filter.RequestedBy).Where("shipments.is_deleted = false").
		Where("shipments.status != ?", constants.ShipmentCreated)

	if filter.ShipmentNature != "" {
		query = query.Where("shipments.shipment_nature = ?", filter.ShipmentNature)
	}

	if filter.Dashboard == constants.Procurement {
		query = query.Where("c.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("shipments.sales_executive_id = ?", filter.RequestedBy)
	}

	if filter.Type != "" {
		if filter.Type == constants.ShipmentTypeFCL+constants.ShipmentTypeLCL {
			query = query.Where("(shipments.type = ? OR shipments.type = ?)", constants.ShipmentTypeFCL, constants.ShipmentTypeLCL)
		} else {
			query = query.Where("shipments.type = ?", filter.Type)
		}
	}

	if filter.Zone != "" {
		if len(ports) > 0 {
			query = query.Where("shipments.pol IN (?) OR shipments.pod IN (?)", ports, ports)
		} else {
			query = query.Where("shipments.id is null")
		}
	}

	if filter.From > 0 {
		query = query.Where("shipments.created_at >= ?", time.Unix(filter.From, 0))
	}

	if filter.To > 0 {
		query = query.Where("shipments.created_at <= ?", time.Unix(filter.To, 0))
	}

	err := query.Debug().Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *Shipment) GetInsightDetails(ctx *context.Context, filter *dtos.DashboardFilters, ports []string) (*dtos.InsightDetailsResponse, error) {
	res := &dtos.InsightDetailsResponse{}
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Select(`count(DISTINCT shipments.id), round(sum(DISTINCT shipments.occupied_cbm),2) as volume, round(sum(DISTINCT shipments.occupied_weight),2) as weight, round(sum(DISTINCT bli.total_sell)::NUMERIC, 2) as revenue, round((sum(DISTINCT bli.total_sell)-sum(DISTINCT bli.total_buy))::NUMERIC, 2) as profit, 
		round(sum(DISTINCT bli.total_buy)::NUMERIC, 2) as total_buy, STRING_AGG(DISTINCT shipments.id::TEXT, ',') AS booking_ids , string_agg(DISTINCT bli.line_item_ids, ',') AS line_item_ids,
		  sum(shipments.teus) as teus`).
		Joins(`JOIN (SELECT li.quote_id, sum(DISTINCT li.buy*li.units*buy_ex.exchange_rate)::NUMERIC as total_buy, 
		sum(DISTINCT li.sell*li.units*sell_ex.exchange_rate)::NUMERIC as total_sell,string_agg(DISTINCT li.id::text, ',') as line_item_ids  
		FROM line_items li JOIN line_item_exchange_rates buy_ex ON li.id = buy_ex.line_item_id AND buy_ex.type = 'buyrate' AND buy_ex.region_id = '`+filter.RegionID+
			`' JOIN line_item_exchange_rates sell_ex ON li.id = sell_ex.line_item_id AND sell_ex.type = 'sellrate' AND sell_ex.region_id ='`+filter.RegionID+`' AND li.region_id ='`+filter.RegionID+`' AND buy_ex.region_id ='`+filter.RegionID+`'WHERE li.sub_type != 'Tax' group by li.quote_id) as bli ON (bli.quote_id = shipments.quote_id)`).
		Joins("JOIN rfqs as r ON (r.id = shipments.rfq_id)").
		Joins("LEFT JOIN cards as c ON ((r.id)::TEXT = c.instance_id AND c.assigned_to = ?)", filter.RequestedBy).Where("shipments.is_deleted = false").
		Where("shipments.status != ?", constants.ShipmentCreated)

	switch filter.TaskName {
	case "day":
		query = query.Select(`TO_CHAR(shipments.created_at::DATE, 'dd MON yyyy') as created_at, count(DISTINCT shipments.id), round(sum(DISTINCT shipments.occupied_cbm),2) as volume, round(sum(DISTINCT shipments.occupied_weight),2) as weight, round(sum(DISTINCT bli.total_sell)::NUMERIC, 2) as revenue, round((sum(DISTINCT bli.total_sell)-sum(DISTINCT bli.total_buy))::NUMERIC, 2) as profit, 
		round(sum(DISTINCT bli.total_buy)::NUMERIC, 2) as total_buy, STRING_AGG(DISTINCT shipments.id::TEXT, ',') AS booking_ids , string_agg(DISTINCT bli.line_item_ids, ',') AS line_item_ids, sum(shipments.teus) as teus`)
		query = query.Group("TO_CHAR(shipments.created_at:: DATE, 'dd MON yyyy')")
	case "week":
		query = query.Select(`To_char(date_trunc('week', shipments.created_at)::DATE, 'dd/mm/yyyy')  || ' - ' || To_char(date_trunc('week', shipments.created_at + '6 days'::interval )::DATE, 'dd/mm/yyyy') as created_at, count(DISTINCT shipments.id), round(sum(DISTINCT shipments.occupied_cbm),2) as volume, round(sum(DISTINCT shipments.occupied_weight),2) as weight, round(sum(DISTINCT bli.total_sell)::NUMERIC, 2) as revenue, round((sum(DISTINCT bli.total_sell)-sum(DISTINCT bli.total_buy))::NUMERIC, 2) as profit, 
		round(sum(DISTINCT bli.total_buy)::NUMERIC, 2) as total_buy, STRING_AGG(DISTINCT shipments.id::TEXT, ',') AS booking_ids , string_agg(DISTINCT bli.line_item_ids, ',') AS line_item_ids, sum(shipments.teus) as teus`)
		query = query.Group("To_char(date_trunc('week', shipments.created_at)::DATE, 'dd/mm/yyyy')  || ' - ' || To_char(date_trunc('week', shipments.created_at + '6 days'::interval )::DATE, 'dd/mm/yyyy')")
	case "month":
		query = query.Select(`to_char(shipments.created_at:: DATE, 'MON yyyy') as created_at, count(DISTINCT shipments.id),round(sum(DISTINCT shipments.occupied_cbm),2) as volume, round(sum(DISTINCT shipments.occupied_weight),2) as weight, round(sum(DISTINCT bli.total_sell)::NUMERIC, 2) as revenue, round((sum(DISTINCT bli.total_sell)-sum(DISTINCT bli.total_buy))::NUMERIC, 2) as profit, 
		round(sum(DISTINCT bli.total_buy)::NUMERIC, 2) as total_buy, STRING_AGG(DISTINCT shipments.id::TEXT, ',') AS booking_ids , string_agg(DISTINCT bli.line_item_ids, ',') AS line_item_ids, sum(shipments.teus) as teus`)
		query = query.Group("to_char(shipments.created_at :: DATE, 'MON yyyy')")
	case "quarter":
		query = query.Select(`to_char(date_trunc('quarter',shipments.created_at)::DATE,'dd/mm/yyyy')  as created_at, count(DISTINCT shipments.id), round(sum(DISTINCT shipments.occupied_cbm),2) as volume, round(sum(DISTINCT shipments.occupied_weight),2) as weight, round(sum(DISTINCT bli.total_sell)::NUMERIC, 2) as revenue, round((sum(DISTINCT bli.total_sell)-sum(DISTINCT bli.total_buy))::NUMERIC, 2) as profit, 
		round(sum(DISTINCT bli.total_buy)::NUMERIC, 2) as total_buy, STRING_AGG(DISTINCT shipments.id::TEXT, ',') AS booking_ids , string_agg(DISTINCT bli.line_item_ids, ',') AS line_item_ids, sum(shipments.teus) as teus`)
		query = query.Group("to_char(date_trunc('quarter',shipments.created_at)::DATE,'dd/mm/yyyy')")
	}

	if filter.CompanyID != "" {
		query = query.Where("shipments.company_id = ?", filter.CompanyID)
	}

	if filter.ShipmentNature != "" {
		query = query.Where("(CASE WHEN (shipments.origin_region_id = ? AND shipments.dest_region_id != ?) THEN 'export' WHEN (shipments.origin_region_id != ? AND shipments.dest_region_id = ?) THEN 'import' ELSE shipments.shipment_nature  END) = ?", filter.RegionID, filter.RegionID, filter.RegionID, filter.RegionID, filter.ShipmentNature)
	}

	if filter.Dashboard == constants.Procurement {
		query = query.Where("c.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("shipments.sales_executive_id = ?", filter.RequestedBy)
	}

	if filter.Type != "" {
		if filter.Type == constants.ShipmentTypeFCL+constants.ShipmentTypeLCL {
			query = query.Where("(shipments.type = ? OR shipments.type = ?)", constants.ShipmentTypeFCL, constants.ShipmentTypeLCL)
		} else {
			query = query.Where("shipments.type = ?", filter.Type)
		}
	}

	if filter.Zone != "" && filter.Zone != constants.AllZones {
		if len(ports) > 0 {
			query = query.Where("shipments.pol IN (?) OR shipments.pod IN (?)", ports, ports)
		} else {
			query = query.Where("shipments.id is null")
		}
	}

	if filter.From > 0 {
		query = query.Where("shipments.created_at >= ?", time.Unix(filter.From, 0))
	}

	if filter.To > 0 {
		query = query.Where("shipments.created_at <= ?", time.Unix(filter.To, 0))
	}

	err := query.Debug().Scan(&res.GetInsightDetails).Error
	if err != nil {
		return nil, err
	}

	res.InsightsCount = int64(len(res.GetInsightDetails))

	return res, nil
}

func (t *Shipment) GetInsightTimeRangeProcurement(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.BookingJobs, error) {
	res := &dtos.BookingJobs{}
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Select(`count(DISTINCT shipments.id), round(sum(DISTINCT shipments.occupied_cbm),2) as total_volume, round(sum(DISTINCT shipments.occupied_weight),2) as total_weight, round(sum(DISTINCT bli.total_sell)::NUMERIC, 2) as total_revenue, round((sum(DISTINCT bli.total_sell)-sum(DISTINCT bli.total_buy))::NUMERIC, 2) as total_profit, 
		round(sum(DISTINCT bli.total_buy)::NUMERIC, 2) as total_buy, STRING_AGG(DISTINCT shipments.id::TEXT, ',') AS booking_ids , string_agg(DISTINCT bli.line_item_ids, ',') AS line_item_ids, sum(shipments.teus) as total_teu`).
		Joins(`JOIN (SELECT li.quote_id, sum(DISTINCT li.buy*li.units*buy_ex.exchange_rate)::NUMERIC as total_buy, 
		sum(DISTINCT li.sell*li.units*sell_ex.exchange_rate)::NUMERIC as total_sell,string_agg(DISTINCT li.id::text, ',') as line_item_ids  
		FROM line_items li JOIN line_item_exchange_rates buy_ex ON li.id = buy_ex.line_item_id AND buy_ex.type = 'buyrate' AND buy_ex.region_id = '`+filter.RegionID+
			`' JOIN line_item_exchange_rates sell_ex ON li.id = sell_ex.line_item_id AND sell_ex.type = 'sellrate' AND sell_ex.region_id ='`+filter.RegionID+`' AND li.region_id ='`+filter.RegionID+`' AND buy_ex.region_id ='`+filter.RegionID+`'WHERE li.sub_type != 'Tax' group by li.quote_id) as bli ON (bli.quote_id = shipments.quote_id)`).
		Joins("JOIN rfqs as r ON (r.id = shipments.rfq_id)").
		Joins("LEFT JOIN cards as c ON ((r.id)::TEXT = c.instance_id  AND c.assigned_to = ?)", filter.RequestedBy).Where("shipments.is_deleted = false").
		Where("shipments.status != ?", constants.ShipmentCreated)

	if filter.Dashboard == constants.Procurement {
		query = query.Where("c.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("shipments.sales_executive_id = ?", filter.RequestedBy)
	}

	if filter.Type != "" {
		if filter.Type == constants.ShipmentTypeFCL+constants.ShipmentTypeLCL {
			query = query.Where("(shipments.type = ? OR shipments.type = ?)", constants.ShipmentTypeFCL, constants.ShipmentTypeLCL)
		} else {
			query = query.Where("shipments.type = ?", filter.Type)
		}
	}

	if filter.From > 0 {
		query = query.Where("shipments.created_at >= ?", time.Unix(filter.From, 0))
	}

	if filter.To > 0 {
		query = query.Where("shipments.created_at <= ?", time.Unix(filter.To, 0))
	}

	err := query.Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *Shipment) GetTimeRangeForInsight(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.BookingJobs, error) {
	res := &dtos.BookingJobs{}
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Select(`count(DISTINCT shipments.id), round(sum(DISTINCT shipments.occupied_cbm),2) as total_volume, round(sum(DISTINCT shipments.occupied_weight),2) as total_weight, round(sum(DISTINCT bli.total_sell)::NUMERIC, 2) as total_revenue, round((sum(DISTINCT bli.total_sell)-sum(DISTINCT bli.total_buy))::NUMERIC, 2) as total_profit, 
		round(sum(DISTINCT bli.total_buy)::NUMERIC, 2) as total_buy, STRING_AGG(DISTINCT shipments.id::TEXT, ',') AS booking_ids , string_agg(DISTINCT bli.line_item_ids, ',') AS line_item_ids, sum(shipments.teus) as total_teu`).
		Joins(`JOIN (SELECT li.quote_id, sum(DISTINCT li.buy*li.units*buy_ex.exchange_rate)::NUMERIC as total_buy, 
		sum(DISTINCT li.sell*li.units*sell_ex.exchange_rate)::NUMERIC as total_sell,string_agg(DISTINCT li.id::text, ',') as line_item_ids  
		FROM line_items li JOIN line_item_exchange_rates buy_ex ON li.id = buy_ex.line_item_id AND buy_ex.type = 'buyrate' AND buy_ex.region_id = '`+filter.RegionID+
			`' JOIN line_item_exchange_rates sell_ex ON li.id = sell_ex.line_item_id AND sell_ex.type = 'sellrate' AND sell_ex.region_id ='`+filter.RegionID+`' AND li.region_id ='`+filter.RegionID+`' AND buy_ex.region_id ='`+filter.RegionID+`'WHERE li.sub_type != 'Tax' group by li.quote_id) as bli ON (bli.quote_id = shipments.quote_id)`).
		Joins("JOIN rfqs as r ON (r.id = shipments.rfq_id)").
		Joins("LEFT JOIN cards as c ON ((r.id)::TEXT = c.instance_id  AND c.assigned_to = ?)", filter.RequestedBy).Where("shipments.is_deleted = false").
		Where("shipments.status != ?", constants.ShipmentCreated)

	if filter.Dashboard == constants.Procurement {
		query = query.Where("c.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("shipments.sales_executive_id = ?", filter.RequestedBy)
	}

	if filter.Type != "" {
		if filter.Type == constants.ShipmentTypeFCL+constants.ShipmentTypeLCL {
			query = query.Where("(shipments.type = ? OR shipments.type = ?)", constants.ShipmentTypeFCL, constants.ShipmentTypeLCL)
		} else {
			query = query.Where("shipments.type = ?", filter.Type)
		}
	}

	if filter.From > 0 {
		query = query.Where("shipments.created_at >= ?", time.Unix(filter.From, 0))
	}

	if filter.To > 0 {
		query = query.Where("shipments.created_at <= ?", time.Unix(filter.To, 0))
	}

	err := query.Scan(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *Shipment) GetforFunnelFilterBookings(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.LiveViewResponse, error) {
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select(`shipments.id ,shipments.code, extract(epoch from shipments.created_at)::INTEGER as created_at,
		shipments.pol_name as origin_port, shipments.pod_name as dest_port, shipments.quote_id, shipments.company_name,
		(CASE WHEN shipments.type != 'FCL' THEN round(shipments.occupied_cbm,2) ELSE 0 END) as volume, 
		(CASE WHEN shipments.type != 'FCL' THEN round(shipments.occupied_weight,2) ELSE 0 END) as weight,
		shipments.type, shipments.teus as teu,shipments.company_id ::text AS customer_id,
		(c.completed_at > c.estimate) as is_breached, shipments.shipment_nature
		`).
		Joins("JOIN rfqs as r ON (r.id = shipments.rfq_id)").
		Joins("LEFT JOIN cards as c ON ((r.id)::TEXT = c.instance_id AND c.assigned_to = ?)", filter.RequestedBy).
		Where("shipments.status = ?", constants.ShipmentCreated).Where("shipments.is_deleted = false")

	if filter.Dashboard == constants.Procurement {
		query = query.Where("c.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("(shipments.sales_executive_id = ?)", filter.RequestedBy)
	}

	if filter.From > 0 {
		query = query.Where("shipments.created_at >= ?", time.Unix(filter.From, 0))
	}

	if filter.To > 0 {
		query = query.Where("shipments.created_at <= ?", time.Unix(filter.To, 0))
	}

	if filter.Type == constants.ShipmentTypeFCL || filter.Type == constants.ShipmentTypeLCL || filter.Type == constants.ShipmentTypeAIR {
		if filter.Type == constants.ShipmentTypeFCL+constants.ShipmentTypeLCL {
			query = query.Where("(shipments.type = ? OR shipments.type = ?)", constants.ShipmentTypeFCL, constants.ShipmentTypeLCL)
		} else {
			query = query.Where("shipments.type = ?", filter.Type)
		}
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

func (t *Shipment) GetFunnelViewCount(ctx *context.Context, filter *dtos.DashboardFilters) (int, error) {
	var cnt int
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).
		Joins("JOIN rfqs as r ON (r.id = shipments.rfq_id)").
		Joins("LEFT JOIN cards as c ON ((r.id)::TEXT = c.instance_id AND c.assigned_to = ?)", filter.RequestedBy).Where("shipments.is_deleted = false")

	switch filter.MetricsFilter {
	case "volume":
		query = query.Select("COALESCE(sum(shipments.occupied_cbm), 0)::INT as volume")
	case "weight":
		query = query.Select("COALESCE(sum(shipments.occupied_weight), 0)::INT as weight")
	case "teus":
		query = query.Select("count(shipments.id)")
	default:
		query = query.Select("count(shipments.id)")
	}

	if filter.Dashboard == constants.Procurement {
		query = query.Where("c.name = ?", constants.CardBuyRate)
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("shipments.sales_executive_id = ?", filter.RequestedBy)
	}

	if filter.Type == constants.ShipmentTypeFCL || filter.Type == constants.ShipmentTypeLCL || filter.Type == constants.ShipmentTypeAIR {
		if filter.Type == constants.ShipmentTypeFCL+constants.ShipmentTypeLCL {
			query = query.Where("(shipments.type = ? OR shipments.type = ?)", constants.ShipmentTypeFCL, constants.ShipmentTypeLCL)
		} else {
			query = query.Where("shipments.type = ?", filter.Type)
		}
	}

	if filter.From > 0 {
		query = query.Where("shipments.created_at >= ?", time.Unix(filter.From, 0))
	}

	if filter.To > 0 {
		query = query.Where("shipments.created_at <= ?", time.Unix(filter.To, 0))
	}

	switch filter.ReportType {
	case constants.FilterBookingCreated:
		query = query.Where("shipments.status = ?", constants.ShipmentCreated)
	case constants.FilterBookingConfirmed:
		query = query.Where("shipments.status != ?", constants.ShipmentCreated)
	}

	err := query.Scan(&cnt).Error
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

func (t *Shipment) GetLiveBookings(ctx *context.Context, filter *dtos.DashboardFilters) (*dtos.LiveViewResponse, error) {
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select(`distinct shipments.id ,shipments.code, extract(epoch from shipments.created_at)::INTEGER as created_at,
		shipments.pol_name as origin_port, shipments.pod_name as dest_port, shipments.quote_id, shipments.company_name,
		(CASE WHEN shipments.type != 'FCL' THEN round(shipments.occupied_cbm,2) ELSE 0 END) as volume, 
		(CASE WHEN shipments.type != 'FCL' THEN round(shipments.occupied_weight,2) ELSE 0 END) as weight,
		shipments.type, shipments.teus as teu,
		(c.completed_at > c.estimate) as is_breached, shipments.shipment_nature, shipments.company_id ::text AS customer_id
		`).
		Joins("JOIN rfqs as r ON (r.id = shipments.rfq_id)").
		Joins("LEFT JOIN cards as c ON ((r.id)::TEXT = c.instance_id AND c.assigned_to = ?)", filter.RequestedBy).
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Buy Rates' GROUP BY instance_id) as bre ON ((r.id)::TEXT = bre.instance_id)").
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Sell Rates'  GROUP BY instance_id) as sre ON ((r.id)::TEXT = sre.instance_id)").
		Where("shipments.is_deleted = false")

	if filter.Dashboard == constants.Procurement {
		query = query.Where("c.name = ?", constants.CardBuyRate)

		if filter.Type == "tat_breached" {
			query = query.Where("bre.expired = true")
		}

		if filter.Type == "within_tat" {
			query = query.Where("(bre.expired = false OR bre.expired IS NULL)")
		}
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("shipments.sales_executive_id = ?", filter.RequestedBy)

		if filter.Type == "tat_breached" {
			query = query.Where("sre.expired = true")
		}

		if filter.Type == "within_tat" {
			query = query.Where("(sre.expired = false OR sre.expired IS NULL)")
		}

	}

	switch filter.Status {
	case constants.LiveBookingYettobeConfirmed:
		query = query.Where("shipments.status = ?", constants.ShipmentCreated)
	case constants.LivebookingConfirmed:
		query = query.Where("shipments.status != ?", constants.ShipmentCreated)
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

func (t *Shipment) GetLiveBookingNotBreached(ctx *context.Context, filter *dtos.DashboardFilters) (int, error) {
	var res int
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("count(shipments.id)").
		Joins("JOIN rfqs as r ON (r.id = shipments.rfq_id)").
		Joins("LEFT JOIN cards as c ON ((r.id)::TEXT = c.instance_id AND  name = 'Buy Rates' AND c.assigned_to = ?)", filter.RequestedBy).
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Buy Rates' GROUP BY instance_id) as bre ON ((r.id)::TEXT = bre.instance_id)").
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Sell Rates'  GROUP BY instance_id) as sre ON ((r.id)::TEXT = sre.instance_id)").
		Where("shipments.is_deleted = false")

	if filter.Dashboard == constants.Procurement {
		query = query.Where("c.name = ?", constants.CardBuyRate).Where("(bre.expired = false or bre.expired is null)")
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("shipments.sales_executive_id = ?", filter.RequestedBy).Where("(sre.expired = false or sre.expired is null)")
	}

	switch filter.Status {
	case constants.StatusCreated:
		query = query.Where("shipments.status = ?", constants.ShipmentCreated)
	case constants.BookingConfirmed:
		query = query.Where("shipments.status != ?", constants.ShipmentCreated)
	}

	err := query.Debug().Scan(&res).Error
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (t *Shipment) GetLiveBookingBreached(ctx *context.Context, filter *dtos.DashboardFilters) (int, error) {
	var res int
	query := ctx.DB.Debug().WithContext(ctx.Request.Context()).Table(t.getTable(ctx)).Select("count(shipments.id)").
		Joins("JOIN rfqs as r ON (r.id = shipments.rfq_id)").
		Joins("JOIN cards as c ON ((r.id)::TEXT = c.instance_id AND  name = 'Buy Rates' AND c.assigned_to = ?)", filter.RequestedBy).
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Buy Rates' GROUP BY instance_id) as bre ON ((r.id)::TEXT = bre.instance_id)").
		Joins("LEFT JOIN (select instance_id, bool_or(CASE WHEN (completed_at > estimate OR (estimate < now() AND status != 'Completed')) THEN true ELSE false END) as expired from cards WHERE name = 'Sell Rates'  GROUP BY instance_id) as sre ON ((r.id)::TEXT = sre.instance_id)").
		Where("shipments.is_deleted = false")

	if filter.Dashboard == constants.Procurement {
		query = query.Where("c.name = ?", constants.CardBuyRate).Where("bre.expired = true")
	}

	if filter.Dashboard == constants.Sales {
		query = query.Where("shipments.sales_executive_id = ?", filter.RequestedBy).Where("sre.expired = true")
	}

	switch filter.Status {
	case constants.StatusCreated:
		query = query.Where("shipments.status = ?", constants.ShipmentCreated)
	case constants.BookingConfirmed:
		query = query.Where("shipments.status != ?", constants.ShipmentCreated)
	}

	err := query.Debug().Scan(&res).Error
	if err != nil {
		return 0, err
	}

	return res, nil
}
