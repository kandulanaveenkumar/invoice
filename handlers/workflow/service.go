package workflow

import (
	"net/http"
	"strconv"

	dtos "bitbucket.org/radarventures/forwarder-adapters/dtos/shipments"
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"bitbucket.org/radarventures/forwarder-shipments/services/workflow"
	"bitbucket.org/radarventures/forwarder-shipments/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateService(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	req := &dtos.Services{}
	err := c.BindJSON(&req)
	if err != nil {
		ctx.Log.Error("Unable to bing json", zap.Error(err))
	}

	err = workflow.New().CreateService(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK,
		utils.GetResponse(http.StatusOK, "", "created"),
	)

}

func UpdateService(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	req := &dtos.Services{}
	err := c.BindJSON(&req)
	if err != nil {
		ctx.Log.Error("Unable to bing json", zap.Error(err))
	}

	err = workflow.New().UpdateService(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK,
		utils.GetResponse(http.StatusOK, "", "updated"),
	)
}

func GetService(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	service, err := workflow.New().GetService(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, service)
}

func GetServices(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	pg, err := strconv.Atoi(c.Query("pg"))
	if err != nil {
		c.JSON(http.StatusOK,
			utils.GetResponse(http.StatusBadRequest, "", err.Error()),
		)
		return
	}

	logAndGetContext(ctx)
	services, err := workflow.New().GetServices(ctx, &dtos.GetFilter{
		Q:        c.Query("q"),
		Pg:       pg,
		RegionID: c.Query("region_id"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, services)
}
