package handlers

import (
	"net/http"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	shipmenttracking "bitbucket.org/radarventures/forwarder-shipments/services/shipment-tracking"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
)

func GetShipmentTracking(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetShipmentTracking")
	res, err := shipmenttracking.NewShipmentTrackingService().GetShipmentTracking(c, c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetContainerTrackingInfo(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetContainerTrackingInfo")
	resp, err := shipmenttracking.NewShipmentTrackingService().GetContainerTrackingInfo(c, c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}
