package handlers

import (
	"net/http"
	"strconv"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/services/sis"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/gin-gonic/gin"
)

func GetSISAirInfo(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetSISAirInfo")
	res, err := sis.NewSISService().GetSISAirInfo(c, c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func SendSISAirInfo(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "SendSISAirInfo")
	req := &dtos.SISAirInfoSendReq{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	err := sis.NewSISService().SendSISAirInfo(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": constants.StatusCreated,
	})
}

func SISStatusUpdateInternal(c *context.Context) {
	req := &dtos.SisStatusReq{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	err := sis.NewSISService().UpdateSISShipmentStatus(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": constants.StatusCreated,
	})
}

func GetSISBookingRequestsInfo(c *context.Context) {

	pg, _ := strconv.Atoi(c.Query("pg"))
	req := &dtos.GetSISInfoReq{
		Code: c.Query("code"),
		Pg:   int64(pg),
	}

	res, err := sis.NewSISService().GetSISBookingRequests(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetSISInfo(c *context.Context) {

	req := &dtos.GetReq{}

	c.SetLoggingContext(c.Param("sid"), "GetSISInfo")
	req.Id = c.Param("sid")

	res, err := sis.NewSISService().GetSISInfo(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}

func SendISFSISInfo(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "SendISFSISInfo")
	req := &dtos.SendReq{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	res := &dtos.MessageResponse{}
	var err error

	if req.SisType == "SIS" {
		res, err = sis.NewSISService().SendSISInfo(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
			)
			return
		}
	} else {
		res, err = sis.NewSISService().SendISFInfo(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
			)
			return
		}
	}

	c.JSON(http.StatusOK, res)
}
