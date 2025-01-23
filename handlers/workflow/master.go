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

func CreateWorkflowMaster(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	req := &dtos.WorkflowMasters{}
	err := c.BindJSON(&req)
	if err != nil {
		ctx.Log.Error("Unable to bing json", zap.Error(err))
	}

	err = workflow.New().CreateMaster(ctx, req)
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

func GetMaster(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	master, err := workflow.New().GetMaster(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, master)
}

func GetMasters(c *gin.Context) {
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
	service, err := workflow.New().GetMasters(ctx, &dtos.GetFilter{
		Q:        c.Query("q"),
		Type:     c.Query("type"),
		Pg:       pg,
		RegionID: c.Query("region_id"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, service)
}
