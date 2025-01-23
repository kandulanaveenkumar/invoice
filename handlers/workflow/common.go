package workflow

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-adapters/utils/log"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/constants"

	"bitbucket.org/radarventures/forwarder-adapters/utils/db"
	"github.com/google/uuid"
)

func logAndGetContext(c *context.Context) {
	c.RefID = c.Request.Header.Get("X-Request-ID")

	if c.RefID == "" {
		c.RefID = uuid.New().String()
	}

	cfg := config.Get()

	c.Log = log.New(c.RefID, cfg.AppName, cfg.LogLevel)
	c.DB = db.New()

	c.TenantID = c.GetHeader("tenant_id")
	if cfg.Env == constants.Dev {
		c.TenantID = "public"
	}
}
