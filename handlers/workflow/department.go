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

func CreateDepartment(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	req := &dtos.Departments{}
	err := c.BindJSON(&req)
	if err != nil {
		ctx.Log.Error("Unable to bing json", zap.Error(err))
	}

	err = workflow.New().CreateDepartment(ctx, req)
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

func UpdateDepartment(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	req := &dtos.Departments{}
	err := c.BindJSON(&req)
	if err != nil {
		ctx.Log.Error("Unable to bing json", zap.Error(err))
	}

	err = workflow.New().UpdateDepartment(ctx, req)
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

func GetDepartment(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	logAndGetContext(ctx)
	department, err := workflow.New().GetDepartment(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, department)
}

func GetDepartments(c *gin.Context) {
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
	departments, err := workflow.New().GetDepartments(ctx, &dtos.GetFilter{
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

	c.JSON(http.StatusOK, departments)
}

func GetDepartmentsWithRegionId(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}
	regionId := c.Param("rid")
	logAndGetContext(ctx)
	departments, err := workflow.New().GetDepartmentsWithRegionId(ctx, regionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.GetResponse(http.StatusInternalServerError, "", err.Error()),
		)
		return
	}
	c.JSON(http.StatusOK, departments)
}
