package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AbortStruct struct {
	ctx *gin.Context
}

func AbortWithStatus(c *gin.Context) AbortStruct {
	return AbortStruct{ctx: c}
}

func (a AbortStruct) BadRequest(err error) {
	a.ctx.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

type ErrorResponse struct {
	Error string `json:"error"`
}
