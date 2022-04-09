package helpers

import (
	"github.com/aliykh/reddit-feed/pkg/customErrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response - response for success responses.
type Response struct {
	Data interface{} `json:"data"`
}

func RespondOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func RespondCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

func RespondError(c *gin.Context, err error) {
	data := customErrors.ParseError(err)
	// mb only set description on development env, otherwise do not set it
	// data.Description = err.Error()
	c.JSON(data.ErrStatus, data)
}
