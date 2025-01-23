package handlers

import (
	"encoding/json"
	"net/http"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	billRegeneration "bitbucket.org/radarventures/forwarder-shipments/services/billregeneration"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetBillDetails(c *context.Context) {

	req := dtos.GetBillsReq{}

	queryShipmentId := c.Query("sid")
	queryAwbLinkId := c.Query("awb_link_id")
	consolId := c.Query("consol_id")
	name := c.Params.ByName("name")
	req.AwbLinkId = queryAwbLinkId
	req.ShipmentId = queryShipmentId
	req.ConsolId = consolId
	req.Name = name
	req.UserId = c.Account.ID.String()
	c.SetLoggingContext(queryShipmentId, "GetBillDetails")

	res, err := GetBillDetailsCommon(c, req.Name, &req)
	if err != nil {
		c.Log.Error("error while fetching bills", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, res)

}

func SaveBillDetails(c *context.Context) {

	req := dtos.BillDetailsReq{}
	err := c.BindJSON(&req)
	if err != nil {
		c.Log.Error("Unable to bing json", zap.Error(err))
	}

	queryShipmentId := c.Query("sid")
	name := c.Params.ByName("name")
	req.ShipmentId = queryShipmentId
	req.UserId = c.Account.ID.String()
	req.Name = name
	c.SetLoggingContext(queryShipmentId, "SaveBillDetails")

	res, err := SaveBillDetailsCommon(c, req.Name, &req)
	if err != nil {
		c.Log.Error("error while saving the bills", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, res)

}

func GenerateBillDetails(c *context.Context) {

	req := dtos.GetBillsReq{}

	queryShipmentId := c.Query("sid")
	mode := c.Query("mode")
	name := c.Params.ByName("name")
	querytype := c.Params.ByName("type")
	req.ShipmentId = queryShipmentId
	req.UserId = c.Account.ID.String()
	req.Name = name
	req.Type = querytype
	req.Mode = mode
	c.SetLoggingContext(queryShipmentId, "GenerateBillDetails")

	res, err := GenerateBillDetailsCommon(c, req.Name, &req)
	if err != nil {
		c.Log.Error("error while generating bills", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)

}

func DownloadBillDetails(c *context.Context) {

	req := dtos.DownloadReq{}
	err := c.BindJSON(&req)
	if err != nil {
		c.Log.Error("Unable to bing json", zap.Error(err))
	}

	name := c.Params.ByName("name")
	querytype := c.Params.ByName("type")
	req.UserId = c.Account.ID.String()
	req.Name = name
	req.Type = querytype

	res, err := DownloadBillDetailsCommon(c, req.Name, &req)
	if err != nil {
		c.Log.Error("error while downloading bills", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)

}

// GetAvailableMAWBForLinking
func GetAvailableMAWBForLinking(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetAvailableMAWBForLinking")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	res, err := GenerateMawbLinkCommon(c, shipmentId)
	if err != nil {
		c.Log.Error("error while fetching available master airway bill", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)

}

func UpsertRegenerate(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "UpsertRegenerate")
	bid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid 'bookingId'",
		})
		return
	}

	billId, err := uuid.Parse(c.Param("billId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid 'billId'",
		})
		return
	}

	decoder := json.NewDecoder(c.Request.Body)
	dto_reqs := &dtos.BillsRegenerate{}

	err = decoder.Decode(dto_reqs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Error while decoding request'",
		})
		return
	}

	dto_reqs.ShipmentId = bid
	dto_reqs.BillId = billId
	dto_reqs.BillName = c.Param("billName")

	errs := billRegeneration.NewBillsRegenerateService().UpsertRegenerate(c, dto_reqs)
	if errs != nil {
		c.JSON(http.StatusInternalServerError, errs.Error())
		return
	}

	c.JSON(http.StatusOK, utils.MessageCardRequestSuccessful)
}

func GetRegenerate(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetRegenerate")
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid 'bookingId'",
		})
		return
	}

	billId, err := uuid.Parse(c.Param("billId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid 'billId'",
		})
		return
	}

	dto_reqs := &dtos.BillsRegenerate{}

	dto_reqs.ShipmentId = sid
	dto_reqs.BillId = billId
	dto_reqs.BillName = c.Param("billName")

	dto_res, errResp := billRegeneration.NewBillsRegenerateService().GetRegenerate(c, dto_reqs)

	if errResp != nil {
		c.Log.Error("error while fetching available master airway bill", zap.Error(errResp))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", errResp.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, dto_res)

}
