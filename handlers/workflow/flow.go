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

func CreateWorkflow(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	req := &dtos.Workflows{}
	err := c.BindJSON(&req)
	if err != nil {
		ctx.Log.Error("Unable to bing json", zap.Error(err))
	}

	err = workflow.New().CreateWorkflow(ctx, req)
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

func UpdateWorkflowFlows(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	req := &dtos.Workflows{}
	err := c.BindJSON(&req)
	if err != nil {
		ctx.Log.Error("Unable to bing json", zap.Error(err))
	}

	req.Id = c.Param("id")
	err = workflow.New().UpdateWorkflowFlows(ctx, req)
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

func GetWorkflow(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	workflow, err := workflow.New().GetWorkflow(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func GetWorkflows(c *gin.Context) {
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
	workflows, err := workflow.New().GetWorkflows(ctx, &dtos.GetFilter{
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

	c.JSON(http.StatusOK, workflows)
}

func CreateWorkflowInstance(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	err := workflow.New().CreateWorkflowForInstance(ctx, c.Query("instance_id"), c.Query("instance_type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK,
		utils.GetResponse(http.StatusOK, "", "success"),
	)
}
