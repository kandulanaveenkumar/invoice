package handlers

import (
	"bitbucket.org/radarventures/forwarder-adapters/utils/context"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	ctx := &context.Context{
		Context: c,
	}

	LogAndGetContext(ctx)
	var no int64

	res := ctx.DB.WithContext(ctx.Request.Context()).Raw(`SELECT 1`).Scan(&no)
	if res.Error != nil {
		c.JSON(500, gin.H{
			"message": res.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
