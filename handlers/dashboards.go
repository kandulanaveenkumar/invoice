package handlers

import (
	"net/http"
	"strconv"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/dashboards"
	"go.uber.org/zap"
)

func LiveCounts(c *context.Context) {
	res, err := dashboards.New().LiveCounts(c, &dtos.DashboardFilters{
		Dashboard:   c.Query("dashboard"),
		RequestedBy: c.Query("executive_id"),
		RegionID:    c.Account.RegionID,
	})
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func LiveBookings(c *context.Context) {
	pg, _ := strconv.Atoi(c.Query("pg"))

	res, err := dashboards.New().LiveBookings(c, &dtos.DashboardFilters{
		Dashboard:   c.Query("dashboard"),
		RequestedBy: c.Query("executive_id"),
		Pg:          pg,
		Status:      c.Query("status"),
		Type:        c.Query("type"),
		RegionID:    c.Account.RegionID,
	})
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func LiveEnquiries(c *context.Context) {
	pg, _ := strconv.Atoi(c.Query("pg"))

	res, err := dashboards.New().LiveEnquiries(c, &dtos.DashboardFilters{
		Dashboard:   c.Query("dashboard"),
		RequestedBy: c.Query("executive_id"),
		Pg:          pg,
		Status:      c.Query("status"),
		Type:        c.Query("type"),
		ReportType:  c.Query("quick_action"),
		RegionID:    c.Account.RegionID,
	})
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func LiveQuotes(c *context.Context) {
	pg, _ := strconv.Atoi(c.Query("pg"))

	res, err := dashboards.New().LiveQuotes(c, &dtos.DashboardFilters{
		Dashboard:   c.Query("dashboard"),
		RequestedBy: c.Query("executive_id"),
		Pg:          pg,
		Status:      c.Query("status"),
		Type:        c.Query("type"),
		ReportType:  c.Query("filter"),
		RegionID:    c.Account.RegionID,
	})
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetFunnelViewCounts(c *context.Context) {
	pg, _ := strconv.Atoi(c.Query("pg"))

	from, _ := strconv.ParseInt(c.Query("from"), 10, 64)
	to, _ := strconv.ParseInt(c.Query("to"), 10, 64)

	res, err := dashboards.New().GetFunnelViewCounts(c, &dtos.DashboardFilters{
		Dashboard:     c.Query("dashboard"),
		RequestedBy:   c.Query("executive_id"),
		Pg:            pg,
		From:          from,
		To:            to,
		Status:        c.Query("status"),
		Type:          c.Query("type"),
		ReportType:    c.Query("filter"),
		MetricsFilter: c.Query("metrics_filters"),
		RegionID:      c.Account.RegionID,
	})
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetFunnelDropoffDetails(c *context.Context) {
	pg, _ := strconv.Atoi(c.Query("pg"))

	from, _ := strconv.ParseInt(c.Query("from"), 10, 64)
	to, _ := strconv.ParseInt(c.Query("to"), 10, 64)

	res, err := dashboards.New().GetFunnelDropoffDetails(c, &dtos.DashboardFilters{
		Dashboard:     c.Query("dashboard"),
		RequestedBy:   c.Query("executive_id"),
		Pg:            pg,
		From:          from,
		To:            to,
		Status:        c.Query("status"),
		Type:          c.Query("type"),
		ReportType:    c.Query("filter"),
		MetricsFilter: c.Query("metrics_filters"),
		RegionID:      c.Account.RegionID,
	})
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetTimeRangeForInsight(c *context.Context) {
	pg, _ := strconv.Atoi(c.Query("pg"))

	from, _ := strconv.ParseInt(c.Query("from"), 10, 64)
	to, _ := strconv.ParseInt(c.Query("to"), 10, 64)

	res, err := dashboards.New().GetTimeRangeForInsight(c, &dtos.DashboardFilters{
		Dashboard:     c.Query("dashboard"),
		RequestedBy:   c.Query("executive_id"),
		Pg:            pg,
		From:          from,
		To:            to,
		Status:        c.Query("status"),
		Type:          c.Query("type"),
		ReportType:    c.Query("filter"),
		MetricsFilter: c.Query("metrics_filters"),
		ID:            c.Query("id"),
		Zone:          c.Query("zone"),
		RegionID:      c.Account.RegionID,
	})
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetInsightDetails(c *context.Context) {
	pg, _ := strconv.Atoi(c.Query("pg"))

	from, _ := strconv.ParseInt(c.Query("from"), 10, 64)
	to, _ := strconv.ParseInt(c.Query("to"), 10, 64)

	res, err := dashboards.New().GetInsightDetails(c, &dtos.DashboardFilters{
		Dashboard:      c.Query("dashboard"),
		RequestedBy:    c.Query("executive_id"),
		Pg:             pg,
		From:           from,
		To:             to,
		Status:         c.Query("status"),
		Type:           c.Query("type"),
		ReportType:     c.Query("filter"),
		MetricsFilter:  c.Query("metrics_filters"),
		ID:             c.Query("id"),
		Zone:           c.Query("zone"),
		CompanyID:      c.Query("customer_id"),
		TaskName:       c.Query("task_name"),
		Category:       c.Query("category"),
		ShipmentNature: c.Query("shipment_nature"),
		RegionID:       c.Account.RegionID,
	})
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, res)
}
