package handlers

import (
	"net/http"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/shipmentparty"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/google/uuid"
)

func GetShipmentPartyAddresses(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetShipmentPartyAddresses")
	res, err := shipmentparty.NewShipmentPartyService().GetShipmentPartyAddresses(c, c.Param("sid"), c.Query("party_types"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func SaveContact(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "SaveContact")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	req := &dtos.ContactBook{}
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrJSONDecode),
		)
		return
	}

	err = shipmentparty.NewShipmentPartyService().SaveContact(c, shipmentId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, utils.MessageResourceUpdated)
}

func GetDSRShipmentAddressDetails(c *context.Context) {
	var shipmentids []string
	err := c.BindJSON(&shipmentids)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrJSONDecode),
		)
		return
	}
	res, err := shipmentparty.NewShipmentPartyService().GetShipmentAddressDetailsForDSR(c, shipmentids)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, res)
}
