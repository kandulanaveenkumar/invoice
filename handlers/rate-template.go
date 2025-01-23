package handlers

import (
	"net/http"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	ratetemplate "bitbucket.org/radarventures/forwarder-shipments/services/rate-template"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/google/uuid"
)

func GetRateTemplates(c *context.Context) {

	rateTemplates, err := ratetemplate.NewRateTemplateService().GetRateTemplates(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, rateTemplates)
}

func SaveRateTemplate(c *context.Context) {

	req := &dtos.RateTemplate{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			utils.GetResponse(http.StatusBadRequest, "", utils.ErrJSONDecode),
		)
		return
	}

	if c.Param("id") != "" {
		req.Id, err = uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest,
				utils.GetResponse(http.StatusBadRequest, "", utils.ErrParsingUUID.Error()),
			)
			return
		}
	}

	err = ratetemplate.NewRateTemplateService().SaveRateTemplate(c, req)
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

func DeleteRateTemplate(c *context.Context) {

	err := ratetemplate.NewRateTemplateService().DeleteRateTemplate(c, c.Param("id"))
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

func GetRateTemplate(c *context.Context) {

	rateTemplates, err := ratetemplate.NewRateTemplateService().GetRateTemplate(c, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, rateTemplates)
}
