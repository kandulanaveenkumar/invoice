package handlers

import (
	"errors"
	"net/http"
	"strings"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/database/models"
	"bitbucket.org/radarventures/forwarder-shipments/services/card"
	cardassignment "bitbucket.org/radarventures/forwarder-shipments/services/card-assignment"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func HandleCard(c *context.Context) {

	cardRequest := &dtos.CardRequest{}
	err := c.BindJSON(&cardRequest)
	if err != nil {
		c.Log.Error("unable to bind json", zap.Error(err))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = cardassignment.NewCardAssignmentService().HandleCard(c, cardRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, utils.MessageCardRequestSuccessful)
}

func GetAllCards(c *context.Context) {
	id, _ := uuid.Parse(c.Query("card_id"))
	filter := &models.Card{
		Id:           id,
		AssignedTo:   c.Query("executive_id"),
		CompanyId:    c.Query("company_id"),
		InstanceId:   c.Query("instance_id"),
		InstanceType: c.Query("instance_type"),
		Type:         c.Query("type"),
		Name:         c.Query("name"),
		// MilestoneId:      c.Query("milestone_id"),
		// TaskId:           c.Query("task_id"),
	}

	statusList := []string{}
	status := c.Query("status")
	if status != "" {
		statusList = strings.Split(status, ",")
	}

	res, err := card.NewCardService().GetAllCards(c, filter, statusList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)

}

func GetAssignedTo(c *context.Context) {
	filter := &models.Card{
		InstanceId:   c.Query("instance_id"),
		InstanceType: c.Query("instance_type"),
		Name:         c.Query("name"),
	}
	c.SetLoggingContext(filter.InstanceId, "GetAssignedTo")
	res, err := card.NewCardService().GetAssignedTo(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)

}

func GetExecutiveCard(c *context.Context) {

	id := c.Query("executive_id")
	if id == "" && c.Account != nil {
		id = c.Account.ID.String()
	}

	cardFilter := &dtos.GECardFilter{
		Name:       c.Query("name"),
		CompanyId:  c.Query("company_id"),
		InstanceId: c.Query("instance_id"),
	}
	c.SetLoggingContext(cardFilter.InstanceId, "GetExecutiveCard")

	req := &models.Card{
		AssignedTo: id,
	}

	cardId := c.Query("card_id")

	bookingFilter := &dtos.GEShipmentFilter{
		Code:           c.Query("code"),
		Pol:            c.Query("pol"),
		Pod:            c.Query("pod"),
		ShipmentNature: c.Query("shipment_nature"),
		ShipmentType:   c.Query("shipment_type"),
	}

	if bookingFilter.ShipmentType != "" && strings.EqualFold(bookingFilter.ShipmentType, constants.ShipmentTypeMisc) {
		bookingFilter.ShipmentType = constants.ShipmentTypeMisc
	}

	if c.Query("include_escalations") == "true" {
		req.EscalatedTo = id
		req.EscalatedById = append(req.EscalatedById, id)
	}

	if cardId != "" {
		cardIdUuid, err := uuid.Parse(c.Query("card_id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		req = &models.Card{
			Id: cardIdUuid,
		}
		cardFilter = &dtos.GECardFilter{}
		bookingFilter = &dtos.GEShipmentFilter{}

	}

	res, err := card.NewCardService().GetExecutiveCards(c, req, cardFilter, bookingFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)

}

func GetTeamCardTreeCount(c *context.Context) {

	id := c.Query("executive_id")
	if id == "" && c.Account != nil {
		id = c.Account.ID.String()
	}
	cardFilter := &dtos.GECardFilter{
		Name:       c.Query("name"),
		InstanceId: c.Query("instance_id"),
		CompanyId:  c.Query("company_id"),
	}
	c.SetLoggingContext(cardFilter.InstanceId, "GetTeamCardTreeCount")
	bookingFilter := &dtos.GEShipmentFilter{
		Code:           c.Query("code"),
		Pol:            c.Query("pol"),
		Pod:            c.Query("pod"),
		ShipmentNature: c.Query("shipment_nature"),
		ShipmentType:   c.Query("shipment_type"),
	}
	res, err := card.NewCardService().GetTeamCardCount(c, cardFilter, bookingFilter, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetOrgTreeCount(c *context.Context) {

	id := c.Query("executive_id")
	if id == "" {
		id = c.Account.ID.String()
	}
	god := false
	if c.Query("god") == "true" {
		god = true
	}
	department := c.Query("department")
	regionId := c.GetHeader("req_region_id")
	res, err := card.NewCardService().GetOrgTreeCount(c, id, god, regionId, department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetExecutivesInvolved(c *context.Context) {

	rfqId := c.Query("rfq_id")
	shipmentId := c.Query("shipment_id")
	c.SetLoggingContext(shipmentId, "GetExecutivesInvolved")

	res, err := card.NewCardService().GetExecutivesInvolved(c, rfqId, shipmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func GetCardLookup(c *context.Context) {

	id := c.Query("executive_id")
	if id == "" && c.Account != nil {
		id = c.Account.ID.String()
	}

	filter := &dtos.FilterCardLabel{
		AdminID:    id,
		IsManager:  c.Query("is_manager_view"),
		Querry:     c.Query("q"),
		AssignedTo: c.Query("executive_id"),
	}

	res, err := card.NewCardService().GetCardLookup(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func EscalateAndReassign(c *context.Context) {

	req := &dtos.ReExecCard{}
	err := c.BindJSON(&req)
	if err != nil {
		c.Log.Error("unable to bing json", zap.Error(err))
		return
	}

	aid := c.Query("aid")

	err = card.NewCardService().EscalateAndReassign(c, req, aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "success")
}

func UpdateCardStatusAsCompleted(c *context.Context) {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
		)
		return
	}

	err = card.NewCardService().UpdateCardStatusAsCompleted(c, id.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.MessageCardRequestSuccessful)
}

func GetCardsWithFilters(c *context.Context) {

	cardReq := &dtos.Card{}
	err := c.BindJSON(&cardReq)
	if err != nil {
		c.Log.Error("unable to bind json", zap.Error(err))
	}

	res, err := card.NewCardService().GetCardsWithFilter(c, cardReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)

}

func GetBukCardsByCompany(c *context.Context) {

	if c.Query("company_ids") == "" {
		c.JSON(http.StatusInternalServerError, errors.New("empty company_ids"))
		return
	}

	res, err := card.NewCardService().GetBulkCardsByCompany(c, c.Query("company_ids"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)

}

func GetCardMasterData(c *context.Context) {
	mname := c.Query("mname")
	tname := c.Query("tname")
	res, err := card.NewCardService().GetCardMasterForMilestoneAndtask(c, mname, tname)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}
