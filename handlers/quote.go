package handlers

import (
	"errors"
	"net/http"
	"strconv"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/quote"
	"bitbucket.org/radarventures/forwarder-shipments/services/shipment"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/google/uuid"
)

func CreateQuote(c *context.Context) {
	req := &dtos.Quote{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	req.RfqID = c.Param("id")
	c.SetLoggingContext(c.Param("sid"), "CreateQuote")
	req.CallBackUrl = c.Query("callback_url")
	id, err := quote.NewQuoteService().CreateQuote(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusCreated,
		utils.GetResponse(http.StatusCreated, id, utils.MessageResourceAdded),
	)
}

func UpdateQuote(c *context.Context) {
	req := &dtos.Quote{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	typeBuySell := c.Query("type")
	req.CallBackUrl = c.Query("callback_url")
	req.ID = c.Params.ByName("qid")
	c.SetLoggingContext(req.ID, "UpdateQuote")
	isShipment, err := strconv.ParseBool(c.Query("is_shipment"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	id, shipmentReq, err := quote.NewQuoteService().UpdateQuote(c, req, typeBuySell, isShipment)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	if shipmentReq != nil && !isShipment {
		_, err = shipment.NewShipmentService().CreateConsolShipment(c, shipmentReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
			)
			return
		}
	}

	c.JSON(http.StatusCreated,
		utils.GetResponse(http.StatusCreated, id, utils.MessageResourceUpdated),
	)
}

func PreviewPDF(c *context.Context) {
	req := &dtos.Quote{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	req.RfqID = c.Param("id")
	req.ID = c.Param("qid")
	c.SetLoggingContext(req.RfqID, "PreviewPDF")

	q := quote.NewQuoteService()
	resp, err := q.PreviewQuotePdf(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, resp)

}

func GeneratePDF(c *context.Context) {

	req := &dtos.GenerateQuotePDF{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	rfqId := c.Params.ByName("id")
	c.SetLoggingContext(rfqId, "GeneratePDF")

	q := quote.NewQuoteService()
	resp, err := q.GenerateQuotePDF(c, rfqId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, resp)

}

func ShareQuote(c *context.Context) {

	req := &dtos.QuoteShareReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}
	quoteId := c.Params.ByName("qid")
	c.SetLoggingContext(quoteId, "ShareQuote")

	req.Id = quoteId

	status, err := quote.NewQuoteService().ShareQuote(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, utils.GetResponse(http.StatusOK, quoteId, status))

}

func GetQuotesForListing(c *context.Context) {
	rfqId := c.Params.ByName("id")
	c.SetLoggingContext(rfqId, "GetQuotesForListing")
	if rfqId == "" {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", errors.New("rfq id can't be empty")),
		)
	}
	res, err := quote.NewQuoteService().GetAllQuotesForRfq(c, rfqId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetQuoteWithLineItems(c *context.Context) {

	c.SetLoggingContext(c.Query("id"), "GetQuoteWithLineItems")
	id, err := uuid.Parse(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	quoteId, err := uuid.Parse(c.Param("qid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}
	flowType := c.Query("flow_type")
	isc, err := strconv.ParseBool(c.Query("isc"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	isConsol, err := strconv.ParseBool(c.Query("is_consol"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	res, err := quote.NewQuoteService().GetQuoteWithLineItems(c, id, quoteId, flowType, isc, isConsol)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, res)
}

func ApproveRejectGPForShipmentAndQuote(c *context.Context) {

	quoteId := c.Param("qid")
	c.SetLoggingContext(quoteId, "ApproveRejectGPForShipmentAndQuote")

	if quoteId == "" {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrEmptyQuoteID.Error()),
		)
		return
	}

	req := &dtos.GPApproveRejectReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	req.QuoteID = quoteId

	err := quote.NewQuoteService().HandleGpApproveReject(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, "")
}

func UpdateQuoteConsol(c *context.Context) {
	req := &dtos.Quote{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	typeBuySell := c.Query("type")
	req.ID = c.Params.ByName("qid")
	c.SetLoggingContext(req.ID, "UpdateQuoteConsol")
	id, err := quote.NewQuoteService().UpdateQuoteConsol(c, req, typeBuySell)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusCreated,
		utils.GetResponse(http.StatusCreated, id, utils.MessageResourceUpdated),
	)
}

func GetQuote(c *context.Context) {
	c.SetLoggingContext(c.Param("qid"), "GetQuote")
	quoteId, err := uuid.Parse(c.Param("qid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}
	if quoteId == uuid.Nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	quote, err := quote.NewQuoteService().GetQuote(c, quoteId.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, quote)

}
