package http

import (
	"fmt"
	"github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/http/server/helpers"
	"github.com/aliykh/reddit-feed/internal/posts"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type handlers struct {
	logger *log.Factory
	uc     posts.UseCase
}

func New(logger *log.Factory, uc posts.UseCase) *handlers {
	return &handlers{
		logger: logger,
		uc:     uc,
	}
}

func (h *handlers) Create(c *gin.Context) {

	model := &models.Post{}

	// go-validator validations
	if err := c.ShouldBindJSON(model); err != nil {
		h.logger.Default().Error(fmt.Sprintf("error while binding json body: %v\n", err.Error()))
		helpers.RespondError(c, err)
		return
	}

	// custom validations
	if err := model.CheckValidity(); err != nil {
		h.logger.Default().Error(fmt.Sprintf("post model: %v\n", err.Error()))
		helpers.RespondError(c, err)
		return
	}

	result, err := h.uc.Create(c.Request.Context(), model)

	if err != nil {
		h.logger.Default().Error("post create", zap.String("err", err.Error()))
		helpers.RespondError(c, err)
		return
	}

	helpers.RespondCreated(c, result)
}
