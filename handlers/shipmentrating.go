package handlers

import (
	"net/http"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/shipmentrating"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
)

func AddShipmentRating(c *context.Context) {

	req := &dtos.AddRatingReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	c.SetLoggingContext(c.Param("sid"), "AddShipmentRating")
	req.ShipmentId = c.Params.ByName("sid")

	res, err := shipmentrating.NewShipmentRatingService().AddShipmentRating(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusCreated,
		utils.GetResponse(http.StatusCreated, res.Id, utils.MessageShipmentCreated),
	)
}

func GetShipmentRatingDetails(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetShipmentRatingDetails")
	sid := c.Params.ByName("sid")

	res, err := shipmentrating.NewShipmentRatingService().GetShipmentRatingDetails(c, sid, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetLatestShipmentToRate(c *context.Context) {

	cid := c.Params.ByName("cid")

	res, err := shipmentrating.NewShipmentRatingService().GetLatestShipmentToRate(c, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetAllRatingMastersDetails(c *context.Context) {

	res, err := shipmentrating.NewShipmentRatingService().GetAllRatingMastersDetails(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetCustomerActivities(c *context.Context) {

	cid := c.Params.ByName("cid")

	res, err := shipmentrating.NewShipmentRatingService().GetCustomerActivities(c, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}
