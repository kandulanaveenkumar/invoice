package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"bitbucket.org/radarventures/forwarder-adapters/apis/docs"
	adapterContext "bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-adapters/utils/db"
	ulog "bitbucket.org/radarventures/forwarder-adapters/utils/log"
	"bitbucket.org/radarventures/forwarder-adapters/utils/tracing"
	"bitbucket.org/radarventures/forwarder-shipments/adhocjobs"
	"bitbucket.org/radarventures/forwarder-shipments/config"
	"bitbucket.org/radarventures/forwarder-shipments/constants"
	"bitbucket.org/radarventures/forwarder-shipments/cronjobs"
	"bitbucket.org/radarventures/forwarder-shipments/jobs"
	"bitbucket.org/radarventures/forwarder-shipments/routes"
	"github.com/gin-gonic/gin"
)

const (
	ENVDev = "dev"
)

func main() {
	var file *os.File
	var err error

	env := os.Getenv("ENV")
	if env == "" {
		env = ENVDev
	}

	file, err = os.Open(env + ".json")
	if err != nil {
		log.Println("Unable to open file. Err:", err)
		os.Exit(1)
	}

	var cnf *config.Config
	config.ParseJSON(file, &cnf)
	config.Set(cnf)

	constants.Logger = ulog.New("Global", cnf.AppName, cnf.LogLevel)
	db.Init(&db.Config{
		URL:           cnf.DatabaseURL,
		MaxDBConn:     cnf.MaxDBConn,
		EnableTracing: true,
	})

	config.InitTemplates()

	job := flag.String("job", "", "Flag to check if job need to Run")
	cronjob := flag.String("cronjob", "", "Flag to check if cronjob need to Run")
	flag.Parse()

	docs.InitializeDocs(config.Get().DocsURL, constants.Logger)

	if job != nil && len(*job) > 0 {
		adhocjobs.Start(job)
		os.Exit(0)
	}

	if cronjob != nil && len(*cronjob) > 0 {
		cronjobs.Start(cronjob)
		os.Exit(0)
	}

	if config.Get().EnableTracing {
		ctx := context.Background()
		t := tracing.NewOtlpTracer(os.Getenv("OTEL_SERVICE_NAME"), os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("OTEL_SERVICE_NAMESPACE"))
		shutdown := t.ConfigureOtlp(ctx)
		defer shutdown()
	}

	c := &adapterContext.Context{}
	cfg := config.Get()
	c.Log = ulog.New(c.RefID, cfg.AppName+"-Job", cfg.LogLevel)
	c.DB = db.New()
	c.TenantID = "public"
	c.Context = &gin.Context{
		Request:  &http.Request{},
		Writer:   nil,
		Params:   []gin.Param{},
		Keys:     map[string]any{},
		Errors:   []*gin.Error{},
		Accepted: []string{},
	}

	if env != ENVDev {
		go jobs.ListenToFreightCRMQueue(c)
	}

	r := routes.GetRouter()

	constants.Logger.Info("Listening to Port: " + cnf.Port)
	r.Run(":" + cnf.Port)
}
