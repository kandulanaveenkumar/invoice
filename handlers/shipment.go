package handlers

import (
	"net/http"
	"strconv"
	"strings"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/config/globals"
	"bitbucket.org/radarventures/forwarder-shipments/services/airwaybillinfo"
	"bitbucket.org/radarventures/forwarder-shipments/services/charges"
	"bitbucket.org/radarventures/forwarder-shipments/services/document"
	globalaccounting "bitbucket.org/radarventures/forwarder-shipments/services/global-accounting"
	"bitbucket.org/radarventures/forwarder-shipments/services/quote"
	"bitbucket.org/radarventures/forwarder-shipments/services/shipment"
	"bitbucket.org/radarventures/forwarder-shipments/services/shipmentcontainer"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrEmptyAccountID = "empty or invalid account id"
)

func CreateShipment(c *context.Context) {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadGateway,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}
	c.SetLoggingContext(c.Param("id"), "CreateShipment")

	qid, err := uuid.Parse(c.Param("qid"))
	if err != nil {
		c.JSON(http.StatusBadGateway,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	req := &dtos.CreateShipmentReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}
	req.RfqId = id
	req.QuoteId = qid

	res, err := shipment.NewShipmentService().CreateShipment(c, req)
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

func CreateMiscShipment(c *context.Context) {

	req := &dtos.MiscShipment{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	res, err := shipment.NewShipmentService().CreateMiscShipment(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusCreated,
		utils.GetResponse(http.StatusCreated, res, utils.MessageShipmentCreated),
	)
}

func UpdateMiscShipment(c *context.Context) {

	req := &dtos.MiscShipment{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	res, err := shipment.NewShipmentService().UpdateMiscShipment(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK,
		utils.GetResponse(http.StatusOK, res, utils.MessageShipmentUpdated),
	)
}

func GetShipmentInternal(c *context.Context) {

	Id := c.Params.ByName("sid")
	c.SetLoggingContext(Id, "GetShipmentInternal")
	shipmentId, err := uuid.Parse(Id)
	if err != nil {
		c.Log.Error("unable to parse uuid", zap.Error(err))
		return
	}
	res, err := shipment.NewShipmentService().GetShipmentInternal(c, shipmentId)
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetShipment(c *context.Context) {

	Id := c.Params.ByName("sid")
	c.SetLoggingContext(Id, "GetShipment")
	shipmentId, err := uuid.Parse(Id)
	if err != nil {
		c.Log.Error("unable to parse uuid", zap.Error(err))
		return
	}
	res, err := shipment.NewShipmentService().GetShipment(c, shipmentId)
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetShipmentsListing(c *context.Context) {

	res, err := shipment.NewShipmentService().GetShipmentsPaginated(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetShipmentsCounts(c *context.Context) {

	res, err := shipment.NewShipmentService().GetShipmentsCounts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetShipmentListSearch(c *context.Context) {

	res, err := shipment.NewShipmentService().GetShipmentListSearch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetConsolListing(c *context.Context) {

	pg, _ := strconv.Atoi(c.Query("pg"))
	req := &dtos.ConsolGetReq{
		AdminId:  c.Account.ID.String(),
		Q:        c.Query("q"),
		Id:       c.Query("id"),
		Status:   c.Query("status"),
		Pg:       int64(pg),
		Type:     globals.BookingTypeCONSOL,
		RegionId: c.Account.RegionID,
		Pol:      c.Query("pol"),
	}

	res, err := shipment.NewShipmentService().GetConsolListing(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetShipmentsForCustomer(c *context.Context) {
	cid := c.Query("cid")
	res, err := shipment.NewShipmentService().GetShipmentsForCustomer(c, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetDocuments(c *context.Context) {

	res, err := document.NewDocumentService().GetDocuments(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetGlobalAccounting(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetGlobalAccounting")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	iscurrencyUSD := c.Query("is_currency_usd")
	if iscurrencyUSD == "" {
		iscurrencyUSD = "false"
	}
	isCurrencyUSD, err := strconv.ParseBool(iscurrencyUSD)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingBool.Error()),
		)
		return
	}

	res, err := globalaccounting.NewGlobalAccountingService().GetGlobalAccounting(c, shipmentId, isCurrencyUSD)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func DeleteConsol(c *context.Context) {

	req := &dtos.ConsolDeleteReq{
		AdminId:      c.Account.ID.String(),
		Id:           c.Param("cid"),
		RegionId:     c.Account.RegionID,
		DeleteReason: c.Query("delete_reason"),
	}

	err := shipment.NewShipmentService().DeleteConsol(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, utils.GetResponse(http.StatusOK, "", "consol deleted"))
}

func EditShipment(c *context.Context) {
	req := &dtos.Shipment{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}
	id := c.Params.ByName("sid")
	c.SetLoggingContext(id, "EditShipment")
	shipmentId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	req.Id = shipmentId
	res, err := shipment.NewShipmentService().EditShipment(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetConsolMatchedShipments(c *context.Context) {
	req := &dtos.ConsolGetReq{
		AdminId:  c.Account.ID.String(),
		Id:       c.Param("cid"),
		RegionId: c.Account.RegionID,
	}

	res, err := shipment.NewShipmentService().GetConsolMatchedShipments(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetConsol(c *context.Context) {

	id := c.Param("cid")
	consolId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	res, err := shipment.NewShipmentService().GetConsol(c, consolId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func ConsolShipmentLinkUnlink(c *context.Context) {

	req := &dtos.ConsolLinkUnlinkReq{}
	if err := c.BindJSON(&req); err != nil {
		c.Log.Error("unable to decode the request", zap.Any("", err))
		return
	}

	req.Action = c.Query("action")

	shipmentServ := shipment.NewShipmentService()
	if req.Action == globals.Link {
		res, err := shipmentServ.ConsolBookingLink(c, req)
		if err != nil {
			c.Log.Error("error while linking shipment", zap.Error(err))
			return
		}
		c.JSON(http.StatusOK, res)
		return
	} else if req.Action == globals.Unlink {
		res, err := shipmentServ.ConsolBookingUnlink(c, req)
		if err != nil {
			c.Log.Error("error while unlinking shipment", zap.Error(err))
			return
		}
		c.JSON(http.StatusOK, res)
		return
	} else if req.Action == globals.Shift {
		res, err := shipmentServ.ConsolBookingShift(c, req)
		if err != nil {
			c.Log.Error("error while shifting shipment", zap.Error(err))
			return
		}
		c.JSON(http.StatusOK, res)
		return
	} else {
		return
	}
}

func UpsertConsol(c *context.Context) {

	req := &dtos.ConsolShipment{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	res, err := shipment.NewShipmentService().UpsertConsol(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func GetShipmentCharges(c *context.Context) {

	id := c.Param("sid")
	c.SetLoggingContext(id, "GetShipmentCharges")
	bookingType := c.Query("bookingtype")
	qtype := c.Query("q_type")
	filter := c.Query("filter")

	shipId, err := uuid.Parse(id)
	if bookingType == "" {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Shipment Id",
			})
			return
		}
	}
	typeBuySell := c.Query("type")
	ch := charges.New()
	res, err := ch.FilterChargesForShipment(c, c.Query("q"), c.Query("header"), shipId, bookingType, qtype, filter, typeBuySell)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

func UpsertConsolInfo(c *context.Context) {

	consolId, err := uuid.Parse(c.Param("cid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Consol Id",
		})
		return
	}
	req := &dtos.ConsolInfoReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	res, err := shipment.NewShipmentService().UpsertConsolInfo(c, req, consolId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func GetMpbCompanies(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetShipmentCharges")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	res, err := shipment.NewShipmentService().GetMpbCompanies(c, shipmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"companies": res,
	})
}

func DeleteMpbCompany(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "DeleteMpbCompany")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	companyId, err := uuid.Parse(c.Param("cid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	err = shipment.NewShipmentService().DeleteMpbCompany(c, shipmentId, companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK,
		utils.GetResponse(http.StatusOK, "", utils.MessageResourceDeleted),
	)
}

func AddMpbAddon(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "AddMpbAddon")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	companyId, err := uuid.Parse(c.Query("cid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	var req []*dtos.LineItem
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	err = shipment.NewShipmentService().AddMpbAddon(c, shipmentId, companyId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK,
		utils.GetResponse(http.StatusOK, "", utils.MessageResourceAdded),
	)
}

func GetConsolJobAccounting(c *context.Context) {

	id := c.Param("cid")
	consolId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	res, err := shipment.NewShipmentService().GetConsolJobAccounting(c, consolId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func UpsertDocument(c *context.Context) {

	req := &dtos.Document{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	instanceId, err := uuid.Parse(c.Param("instance_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	req.InstanceId = instanceId

	err = document.NewDocumentService().UpsertDocument(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, "success")
}

func UpdateETAETD(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "UpdateETAETD")
	req := &dtos.Quote{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	quoteID := c.Query("qid")

	req.ID = quoteID
	err := shipment.NewShipmentService().UpdateETAETD(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, "success")
}

func DeleteBooking(c *context.Context) {

	if c.Account == nil || c.Account.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, ErrEmptyAccountID)
		return
	}

	c.SetLoggingContext(c.Param("sid"), "DeleteBooking")
	req := dtos.ConsolDeleteReq{
		AdminId:      c.Account.ID.String(),
		Id:           c.Param("sid"),
		Type:         c.Query("type"),
		DeleteReason: c.Query("delete_reason"),
	}

	err := shipment.NewShipmentService().DeleteBooking(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "success")

}

func GetCustomersActiveShipmentCount(c *context.Context) {

	reqIds := c.Query("company_ids")
	cids := strings.Split(reqIds, ",")
	res, err := shipment.NewShipmentService().CustomerPaginationAciveBookingsCount(c, cids)
	if err != nil {
		c.Log.Error("error while fetching shipment", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func UpdateBookingStatusAdmin(c *context.Context) {

	if c.Account == nil || c.Account.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, ErrEmptyAccountID)
		return
	}

	c.SetLoggingContext(c.Param("sid"), "UpdateBookingStatusAdmin")
	req := &dtos.StatusUpdateReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := shipment.NewShipmentService().UpdateBookingStatus(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "success")

}

func AddPartnerInvoices(c *context.Context) {

	if c.Account == nil || c.Account.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, ErrEmptyAccountID)
		return
	}

	req := &dtos.PartnerInvoices{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if c.Param("sid") != "" && req.InstanceId == "" {
		req.InstanceId = c.Param("sid")
	}
	c.SetLoggingContext(c.Param("sid"), "AddPartnerInvoices")

	res, err := shipment.NewShipmentService().AddPartnerInvoices(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)

}

func ApproveDeviatedPartnerInvoice(c *context.Context) {

	if c.Account == nil || c.Account.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, ErrEmptyAccountID)
		return
	}

	req := &dtos.PartnerInvoices{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if c.Param("sid") != "" && req.InstanceId == "" {
		req.InstanceId = c.Param("sid")
	}
	c.SetLoggingContext(c.Param("sid"), "ApproveDeviatedPartnerInvoice")

	err := shipment.NewShipmentService().ApproveDeviatedPartnerInvoice(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, "success")

}

func GetBookingTimeline(c *context.Context) {

	sid := c.Param("sid")
	c.SetLoggingContext(sid, "GetBookingTimeline")
	res, err := shipment.NewShipmentService().GetTimelineForBooking(c, sid)
	if err != nil {
		c.Log.Error("error while fetching shipment timeline", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)

}

func GetDSRShipments(c *context.Context) {
	req := &dtos.DSRShipmentParamters{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrJSONDecode),
		)
		return
	}
	res, err := shipment.NewShipmentService().GetDSRShipments(c, req)
	if err != nil {
		c.Log.Error("error while fetching dsr shipments", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, res)
}
func GetAirwayBillNumbersForDSR(c *context.Context) {
	var shipmentids []string
	err := c.BindJSON(&shipmentids)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrJSONDecode),
		)
		return
	}
	res, err := airwaybillinfo.NewAirwayBillInfoService().GetAirwayBillNumbersForDSR(c, shipmentids)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetTaxPercentage(c *context.Context) {

	partnerCountry := c.Query("partner_country")
	partnerState := c.Query("partner_state")
	orgTaxIdNo := c.Query("org_tax_id_no")
	orgState := c.Query("org_state")
	orgCountry := c.Query("org_country")
	billingInstanceId := c.Query("billing_instances_id")
	c.SetLoggingContext(c.Param("sid"), "GetTaxPercentage")

	req := &dtos.TaxFilters{
		PartnerCountry:     partnerCountry,
		PartnerState:       partnerState,
		ShipmentId:         c.Param("sid"),
		PartnerId:          c.Param("pid"),
		ApproveId:          "",
		RefId:              "",
		OrgTaxIdNo:         orgTaxIdNo,
		OrgState:           orgState,
		OrgCountry:         orgCountry,
		BillingInstancesId: billingInstanceId,
		AdminId:            c.Account.ID.String(),
	}

	res, err := shipment.NewShipmentService().GetTaxPercentage(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}
func GetFreightCertificate(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetFreightCertificate")
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid 'Shipment'",
		})
	}

	res, err := shipment.NewShipmentService().GetFreightCertificate(c, sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func UpsertFreightCertificate(c *context.Context) {

	dto_reqs := &dtos.FreightCertificateReq{}

	if err := c.BindJSON(&dto_reqs); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.SetLoggingContext(c.Param("sid"), "UpsertFreightCertificate")
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid 'shipmentId'",
		})
	}

	err = shipment.NewShipmentService().UpsertFreightCertificate(c, dto_reqs, sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, "success")
}

func GenerateFreightCertificate(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GenerateFreightCertificate")
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid 'shipmentId'",
		})
	}

	res, err := shipment.NewShipmentService().GenerateFreightCertificate(c, sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func UpdateShipmentLock(c *context.Context) {

	shipmentLockReq := &dtos.ShipmentLock{}

	if err := c.BindJSON(&shipmentLockReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := shipment.NewShipmentService().UpdateShipmentLockstatus(c, shipmentLockReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func PreAlert(c *context.Context) {

	sid := c.Param("sid")
	c.SetLoggingContext(sid, "PreAlert")
	req := &dtos.PreAlertReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	res, err := shipment.NewShipmentService().UploadPreAlert(c, sid, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)

}

func UpdateContainer(c *context.Context) {

	shipmentContainers := &dtos.ShipmentContainer{}

	if err := c.BindJSON(&shipmentContainers); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := shipmentcontainer.NewShipmentContainerService().DeleteContianer(c, shipmentContainers)

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func RevertShipment(c *context.Context) {

	c.SetLoggingContext(c.Params.ByName("sid"), "RevertShipment")
	rfqId, err := shipment.NewShipmentService().RevertShipment(c, c.Params.ByName("sid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.GetResponse(http.StatusOK, rfqId, utils.MessageShipmentReverted))

}

func GetPartnerInvoices(c *context.Context) {

	sid := c.Param("sid")
	c.SetLoggingContext(sid, "GetPartnerInvoices")

	res, err := quote.NewQuoteService().GetPartnerInvoices(c, sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func UpateCompanyNameForShipmentsAndRfqs(c *context.Context) {

	req := &dtos.CustomerNameChangeReq{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := shipment.NewShipmentService().UpdateCompanyNameforshipmentsAndRfqs(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, "success")
}

func DeleteDocuments(c *context.Context) {
	instanceId := c.Param("instance_id")
	flowInstanceId := c.Param("flow_instance_id")

	err := document.NewDocumentService().DeletetDocuments(c, instanceId, flowInstanceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "success")

}
