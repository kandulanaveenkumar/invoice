package handlers

import (
	"net/http"
	"strings"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/services/charges"
	"bitbucket.org/radarventures/forwarder-shipments/services/rfq"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateRfq(c *context.Context) {

	req := &dtos.RfqReqRes{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	id, refId, err := rfq.NewRfqService().CreateRfq(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"code":    refId,
		"message": constants.StatusCreated,
	})
}

func GetRfq(c *context.Context) {

	c.SetLoggingContext(c.Param("id"), "GetRfq")
	res, err := rfq.NewRfqService().GetRfq(c, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func UpdateRfq(c *context.Context) {

	c.SetLoggingContext(c.Param("id"), "UpdateRfq")
	req := &dtos.RfqReqRes{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	req.ID = uuid.MustParse(c.Param("id"))
	err := rfq.NewRfqService().UpdateRfq(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK,
		utils.GetResponse(http.StatusOK, req.ID.String(), constants.StatusUpdated),
	)
}

func GetCountForListing(c *context.Context) {

	req := &dtos.GetRfqsFiltersReq{}

	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}
	if req.Pg <= 0 || req.Count <= 0 {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", "page or count can't be zero or negative"),
		)
		return
	}

	if c.Query("req_ids") != "" {
		req.Ids = strings.Split(c.Query("req_ids"), ",")
	}

	res, err := rfq.NewRfqService().GetCountForListing(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetRfqs(c *context.Context) {

	req := &dtos.GetRfqsFiltersReq{}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	if req.Pg <= 0 || req.Count <= 0 {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", "page or count can't be zero or negative"),
		)
		return
	}

	if c.Query("req_ids") != "" {
		req.Ids = strings.Split(c.Query("req_ids"), ",")
	}

	if c.Query("partner_id") != "" {
		req.PartnerID = c.Query("partner_id")

	}

	if c.Query("sort_by") != "" {
		req.SortBy = c.Query("sort_by")
	}

	res, err := rfq.NewRfqService().GetAllRfqsForListing(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}
func GetRfqSearch(c *context.Context) {
	c.RegionId = c.GetHeader("switch_region_id")
	req := &dtos.GetRfqsFiltersReq{
		Q: c.Query("q"),
	}
	res, err := rfq.NewRfqService().GetRfqSearchList(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetRfQDefaultLineItems(c *context.Context) {
	rfqId := c.Param("id")
	if err := uuid.Validate(rfqId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid RFQ Id",
		})
		return
	}
	c.SetLoggingContext(rfqId, "GetRfQDefaultLineItems")

	ch := charges.New()
	res, err := ch.GetRfqDefaultCharges(c, rfqId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"line_items": res,
	})
}

func GetRfQCharges(c *context.Context) {
	id := c.Param("id")
	c.SetLoggingContext(id, "GetRfQCharges")
	rfqId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid RFQ Id",
		})
		return
	}
	typeBuySell := c.Query("type")

	ch := charges.New()
	res, err := ch.FilterChargesForRFQ(c, c.Query("q"), c.Query("header"), rfqId, constants.WorkflowTypeRFQ, typeBuySell)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

func GetConsolDefaultLineItems(c *context.Context) {

	consolId, err := uuid.Parse(c.Query("cid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Consol Id",
		})
		return
	}

	ch := charges.New()
	res, err := ch.GetForCONSOL(c, consolId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"line_items": res,
	})
}

func GetConsolCharges(c *context.Context) {
	id := c.Param("cid")
	consolId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid consol Id",
		})
		return
	}

	ch := charges.New()
	res, err := ch.FilterChargesForConsol(c, c.Query("q"), c.Query("header"), consolId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

func GetMiscBookingDefaultLineItems(c *context.Context) {

	ch := charges.New()
	res, err := ch.GetMiscBookingDefaultLineItems(c, c.Query("booking_type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}
func ExpireRfq(c *context.Context) {

	c.SetLoggingContext(c.Param("id"), "ExpireRfq")
	req := &dtos.RfqExpiryReq{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	req.Id = c.Param("id")

	res, err := rfq.NewRfqService().ExpireRfq(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
