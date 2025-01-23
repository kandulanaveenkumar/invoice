package handlers

import (
	"net/http"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/shipment"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func CustomerDashboardDetails(c *context.Context) {

	Id := c.Params.ByName("cid")
	companyId, err := uuid.Parse(Id)
	if err != nil {
		c.Log.Error("unable to parse uuid", zap.Error(err))
		return
	}
	res, err := shipment.NewShipmentService().CustomerDashboardDetails(c, companyId.String())
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func ShipmentTrackingForCustomerDasboard(c *context.Context) {

	Id := c.Params.ByName("cid")
	companyId, err := uuid.Parse(Id)
	if err != nil {
		c.Log.Error("unable to parse uuid", zap.Error(err))
		return
	}
	res, err := shipment.NewShipmentService().ShipmentTrackingForCustomerDasboard(c, companyId.String())
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, res)
}
