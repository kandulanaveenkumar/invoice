package cronjobs

import (
	"net/http/httptest"
	"os"
	"time"

	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-adapters/utils/db"
	ulog "bitbucket.org/radarventures/forwarder-adapters/utils/log"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/services/rfq"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Start(cronjob *string) {

	if *cronjob == "updateQuoteExpiry" {
		ctx := getContext()
		// To mock request context of gin. The request context is used in DAO layer
		ctx.Context, _ = gin.CreateTestContext(httptest.NewRecorder())
		ctx.Context.Request = httptest.NewRequest("GET", "/expired-quotes-job", nil)
		ctx.Log.Info("Started quote expiry job", zap.String("start_time", time.Now().UTC().String()))
		rfqService := rfq.NewRfqService()
		rfqService.UpdateQuoteExpiry(ctx)
		ctx.Log.Info("Completed quote expiry job", zap.Any("end_time", time.Now().UTC().String()))
		os.Exit(0)
	}

	if *cronjob == "handleCardStatus" {
		ctx := getContext()
		ctx.Context, _ = gin.CreateTestContext(httptest.NewRecorder())
		ctx.Context.Request = httptest.NewRequest("GET", "/handle-card-status", nil)
		NewHandleCardStatus().HandleStatus(ctx)
		os.Exit(0)
	}

	if *cronjob == "containerTracking" {
		ctx := getContext()
		ctx.Context, _ = gin.CreateTestContext(httptest.NewRecorder())
		ctx.Context.Request = httptest.NewRequest("GET", "/container-tracking", nil)
		RunContainerTrackingAutomation(ctx)
		os.Exit(0)
	}

	if *cronjob == "sisrequestjob" {
		ctx := getContext()
		ctx.Context, _ = gin.CreateTestContext(httptest.NewRecorder())
		ctx.Context.Request = httptest.NewRequest("GET", "/container-tracking", nil)
		Run(ctx)
		os.Exit(0)
	}

	if *cronjob == "CardsDeletePostExpiry" {
		ctx := getContext()
		ctx.Context, _ = gin.CreateTestContext(httptest.NewRecorder())
		ctx.Context.Request = httptest.NewRequest("GET", "/cards-delete-post-expiry", nil)
		NewCardsPostExpiry().CardsPostExpiry(ctx)
		os.Exit(0)
	}

	if *cronjob == "UpdateShipmentlockStatus" {
		ctx := getContext()
		ctx.Context, _ = gin.CreateTestContext(httptest.NewRecorder())
		ctx.Context.Request = httptest.NewRequest("GET", "/update-shipment-lock", nil)
		NewUpdateShipmentLock().UpdateShipmentlockStatus(ctx)
		os.Exit(0)
	}
	if *cronjob == "UpdateShipmentlockStatuV2" {
		ctx := getContext()
		ctx.Context, _ = gin.CreateTestContext(httptest.NewRecorder())
		ctx.Context.Request = httptest.NewRequest("GET", "/update-shipment-lock", nil)
		NewUpdateShipmentLock().UpdateShipmentlockStatusV2(ctx)
		os.Exit(0)
	}

	if *cronjob == "AMSCheckMail" {
		ctx := getContext()
		ctx.Context, _ = gin.CreateTestContext(httptest.NewRecorder())
		ctx.Context.Request = httptest.NewRequest("GET", "/check-ams", nil)
		CheckMails(ctx)
		os.Exit(0)
	}

	if *cronjob == "InvoiceRetrievel" {
		ctx := getContext()
		ctx.Context, _ = gin.CreateTestContext(httptest.NewRecorder())
		ctx.Context.Request = httptest.NewRequest("GET", "/invoice-retrievel", nil)
		FetchPendingInvoices(ctx)
		os.Exit(0)
	}

	if *cronjob == "expiryCardsNotifications" {
		ctx := getContext()
		ctx.Context, _ = gin.CreateTestContext(httptest.NewRecorder())
		ctx.Context.Request = httptest.NewRequest("GET", "/expiry-cards-notifications", nil)
		NewExpiryCards().SendExpiryCardsNotifications(ctx)
		os.Exit(0)
	}

}

func getContext() *context.Context {
	c := &context.Context{}

	c.RefID = uuid.New().String()

	cfg := config.Get()
	c.Log = ulog.New(c.RefID, cfg.AppName, cfg.LogLevel)
	c.DB = db.New()

	if c.TenantID == "" {
		c.TenantID = "public"
	}

	return c
}
