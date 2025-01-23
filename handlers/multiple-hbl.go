package handlers

import (
	"net/http"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	multiplehbl "bitbucket.org/radarventures/forwarder-shipments/services/multiple-hbl"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetBlInstructions(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetBlInstructions")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	resp, err := multiplehbl.NewMultipleHblService().GetBlInstructions(c, shipmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func SaveBlInstructions(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "SaveBlInstructions")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	req := &dtos.BlInstruction{}
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrJSONDecode.Error()),
		)
		return
	}

	err = multiplehbl.NewMultipleHblService().SaveBlInstructions(c, shipmentId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK,
		utils.GetResponse(http.StatusOK, "", utils.MessageResourceUpdated),
	)
}

func SaveMasterContainers(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "SaveMasterContainers")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadGateway, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	req := &dtos.ContainersList{}
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrJSONDecode.Error()),
		)
		return
	}

	err = multiplehbl.NewMultipleHblService().SaveMasterContainers(c, shipmentId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK,
		utils.GetResponse(http.StatusOK, "", utils.MessageResourceUpdated),
	)
}

func GetContainerMapping(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetContainerMapping")
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	resp, err := multiplehbl.NewMultipleHblService().GetContainerMapping(c, sid)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func CreateBlNo(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "CreateBlNo")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	consigneeId, err := uuid.Parse(c.Query("consigID"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	shipperId, err := uuid.Parse(c.Query("shiperID"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	notifyPartyId, err := uuid.Parse(c.Query("notifyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	resp, err := multiplehbl.NewMultipleHblService().CreateBlNo(c, shipmentId, consigneeId, shipperId, notifyPartyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetBlData(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetBlData")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	resp, err := multiplehbl.NewMultipleHblService().GetBlData(c, shipmentId, c.Param("blno"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func SaveMarksAndDescription(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "SaveMarksAndDescription")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	containerId, err := uuid.Parse(c.Param("conID"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	req := &dtos.MarksAndDescriptions{}
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrJSONDecode.Error()),
		)
		return
	}

	resp, err := multiplehbl.NewMultipleHblService().SaveMarksAndDescription(c, shipmentId, c.Param("blno"), containerId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetHbl(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetHbl")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	resp, err := multiplehbl.NewMultipleHblService().GetHbl(c, shipmentId, c.Param("blno"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func CreateHbl(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "CreateHbl")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	resp, err := multiplehbl.NewMultipleHblService().CreateHbl(c, shipmentId, c.Param("blno"), c.Param("type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func DownloadHbl(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "DownloadHbl")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	docType := c.Param("type")
	if docType != constants.BlTypeOriginal && docType != constants.BlTypeCopy && docType != constants.BlTypeDraft {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid HBL Type  request",
		})
	}

	resp, err := multiplehbl.NewMultipleHblService().DownloadHbl(c, shipmentId, c.Param("blno"), c.Param("type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetGeneratedHblList(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetGeneratedHblList")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	resp, err := multiplehbl.NewMultipleHblService().GetGeneratedHblList(c, shipmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func DeleteBlNo(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "DeleteBlNo")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	resp, err := multiplehbl.NewMultipleHblService().DeleteBlNo(c, shipmentId, c.Param("blno"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func DeleteMarksAndDescriptionsContainer(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "DeleteMarksAndDescriptionsContainer")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	containerID, err := uuid.Parse(c.Param("conID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid User Id request",
		})
	}

	resp, err := multiplehbl.NewMultipleHblService().DeleteMarksAndDescriptionsContainer(c, shipmentId, c.Param("blno"), containerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func DeleteMarksAndDescription(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "DeleteMarksAndDescription")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	containerID, err := uuid.Parse(c.Param("conID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid User Id request",
		})
	}

	marksAndDescriptionId, err := uuid.Parse(c.Param("mid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid marks and desc  Id request",
		})
	}

	resp, err := multiplehbl.NewMultipleHblService().DeleteMarksAndDescription(c, shipmentId, c.Param("blno"), containerID, marksAndDescriptionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetMultipleHbl(c *context.Context) {

	c.SetLoggingContext(c.Param("sid"), "GetMultipleHbl")
	shipmentId, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	resp, err := multiplehbl.NewMultipleHblService().GetMultipleHbl(c, shipmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}
