package handlers

import (
	"net/http"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/services/accountpreference"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
)

func SetAccountPreferences(c *context.Context) {

	req := &dtos.AccountPreference{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	err := accountpreference.NewAccountPreferenceService().SetAccountPreferences(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, constants.StatusCreated)
}

func GetAccountPreference(c *context.Context) {

	res, err := accountpreference.NewAccountPreferenceService().GetAccountPreference(c, c.Query("account_id"), c.Query("type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, res)
}
