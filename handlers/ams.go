package handlers

import (
	"net/http"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/ams"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateAMS(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GenerateAMS")
	req := &dtos.AmsInfo{}

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrJSONDecode.Error()),
		)
		return
	}

	msgs, err := ams.New().GenerateAMS(c, req)
	if err != nil {
		if len(msgs) > 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"errs": msgs,
			})
			return
		}
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, utils.MessageAmsGenerated)
}

func GetAmsInfo(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetAmsInfo")
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	ams, err := ams.New().GetAMSInfo(c, sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, ams)
}

func CompleteRefileAms(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "CompleteRefileAms")
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	hblNo := c.Param("hbl_no")
	if hblNo == "" {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrEmptyBlNo.Error()),
		)
		return
	}

	err = ams.New().CompleteRefileAms(c, sid.String(), hblNo, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, utils.MessageRefileAmsCompleted)
}
