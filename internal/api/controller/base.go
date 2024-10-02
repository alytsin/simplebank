package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Base struct {
}

func (c *Base) validateJsonOrSendBadRequest(ctx *gin.Context, v interface{}) bool {
	if err := ctx.ShouldBindJSON(v); err != nil {
		c.doBadRequest(ctx, err)
		return false
	}
	return true
}

func (c *Base) validateQueryOrSendBadRequest(ctx *gin.Context, v interface{}) bool {
	if err := ctx.ShouldBindQuery(v); err != nil {
		c.doBadRequest(ctx, err)
		return false
	}
	return true
}

func (c *Base) doBadRequest(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, ErrorMessage{Error: err.Error()})
}

func (c *Base) validateUriOrSendBadRequest(ctx *gin.Context, v interface{}) bool {
	if err := ctx.ShouldBindUri(v); err != nil {
		c.doBadRequest(ctx, err)
		return false
	}
	return true
}
