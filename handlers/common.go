package handlers

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-adapters/utils/log"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/constants"

	"bitbucket.org/radarventures/forwarder-adapters/utils/db"
	"github.com/google/uuid"
)

func LogAndGetContext(c *context.Context) {
	c.RefID = c.Request.Header.Get("X-Request-Shipment")

	if c.RefID == "" {
		c.RefID = uuid.New().String()
	}

	cfg := config.Get()

	c.Log = log.New(c.RefID, cfg.AppName, cfg.LogLevel)
	c.DB = db.New()

	c.TenantID = c.GetHeader("tenant_id")
	// Local-Testing:
	if cfg.Env == constants.Dev {
		c.TenantID = "public"
	}

	c.RegionId = c.GetHeader("switch_region_id")
}
