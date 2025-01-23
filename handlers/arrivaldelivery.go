package handlers

import (
	"net/http"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/arrivaldelivery"
	"github.com/gin-gonic/gin"
)

func GetArrivalDeliveryDetails(c *context.Context) {

	req := &dtos.GetBillsReq{}

	c.SetLoggingContext(c.Param("sid"), "GetArrivalDeliveryDetails")
	shipmentId := c.Param("sid")
	blNo := c.Query("bl_no")

	req.ShipmentId = shipmentId
	req.UserId = c.Account.ID.String()
	req.BlNo = blNo

	result, err := arrivaldelivery.NewArrivalDeliveryService().GetArrivalDeliveryDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func SaveArrivalDeliveryDetails(c *context.Context) {

	req := &dtos.ArrivalDeliveryDocs{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.SetLoggingContext(c.Param("sid"), "SaveArrivalDeliveryDetails")
	req.ShipmentId = c.Param("sid")
	req.UserId = c.Account.ID.String()

	result, err := arrivaldelivery.NewArrivalDeliveryService().SaveArrivalDeliveryDetails(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
