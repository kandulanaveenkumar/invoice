package handlers

import (
	"net/http"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	fcr "bitbucket.org/radarventures/forwarder-shipments/services/fcr"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/google/uuid"
)

func GetFcr(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetFcr")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	resp, err := fcr.NewFcrService().GetFcr(c, shipmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func SaveFcr(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "SaveFcr")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	req := &dtos.FcrFields{}
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrJSONDecode.Error()),
		)
		return
	}

	resp, err := fcr.NewFcrService().SaveFcr(c, req, shipmentId, c.Query("blno"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GenerateFcr(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GenerateFcr")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	blNo := c.Query("blno")

	resp, err := fcr.NewFcrService().GenerateFcr(c, shipmentId, blNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)

}
