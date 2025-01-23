package handlers

import (
	"net/http"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/airwaybillmanifest"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetManifest(c *context.Context) {

	req := &dtos.GetBillsReq{}

	shipmentId := c.Query("sid")
	c.SetLoggingContext(shipmentId,"GetManifest")
	req.ShipmentId = shipmentId
	req.UserId = c.Account.ID.String()

	result, err := airwaybillmanifest.NewAirwayBillManifestService().GetAirwayBillManifest(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func CreateManifest(c *context.Context) {

	req := dtos.AirwayBillManifest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.Log.Error("error while binding json", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	shipmentId := c.Query("sid")
	c.SetLoggingContext(shipmentId,"CreateManifest")
	req.ShipmentId = shipmentId
	req.UserId = c.Account.ID.String()

	result, err := airwaybillmanifest.NewAirwayBillManifestService().SaveAirwayBillManifest(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)

}
