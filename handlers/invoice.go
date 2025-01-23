package handlers

import (
	"net/http"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/invoice"
	"bitbucket.org/radarventures/forwarder-shipments/services/invoice/invoicepref"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateInvoicePref(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "CreateInvoicePref")
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadGateway,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	err = invoicepref.NewInvoicePrefsService().CreateInvoicePref(c, sid, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, utils.MessageInvoicePrefCreated)
}

func GetInvoice(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetInvoice")
	res, err := invoice.NewInvoiceService().GetInvoice(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetInvoiceListing(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetInvoiceListing")
	res, err := invoice.NewInvoiceService().GetInvoiceListing(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func GenerateInvoice(c *context.Context) {

	req := &dtos.InvoiceGenerate{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := invoice.NewInvoiceService().GenerateInvoice(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func SaveInvoicePref(c *context.Context) {

	req := &dtos.InvoicePref{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Parse region ID from the query or the account context
	if c.Query("rid") != "" {
		regionId, err := uuid.Parse(c.Query("rid"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		req.RegionId = regionId
	}

	err := invoicepref.NewInvoicePrefsService().SaveInvoicePref(c, req, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.MessageInvoicePrefUpdated)
}

func ShareInvoice(c *context.Context) {

	req := &dtos.InvoiceShare{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := invoice.NewInvoiceService().ShareInvoice(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.MessageInvoiceShared)

}

func GetInvoiceRetrieval(c *context.Context) {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadGateway,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	res, err := invoice.NewInvoiceService().GetInvoiceRetrieval(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func GenerateIcaInvoice(c *context.Context) {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadGateway,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	res, err := invoice.NewInvoiceService().GenerateIcaInvoice(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
