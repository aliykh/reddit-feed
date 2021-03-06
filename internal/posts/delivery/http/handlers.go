package http

import (
	"fmt"
	"github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/http/server/helpers"
	"github.com/aliykh/reddit-feed/internal/posts"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/aliykh/reddit-feed/pkg/pagination"
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

// Create godoc
// @Summary Create - create a new post
// @Description - create a new post
// @Tags Posts
// @Param params body models.Post true "body"
// @Accept json
// @Produce json
// @Success 200 {object} models.Post
// @Router /post [POST]
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

// Generate godoc
// @Summary Generate - generates a feed of posts
// @Description returns a list of posts
// @Tags Posts
// @Accept json
// @Produce json
// @Success 200 {object} models.Feed
// @Router /post/generate [GET]
func (h *handlers) Generate(c *gin.Context) {

	pg := &pagination.Query{
		Size: 25,
	}

	if err := c.ShouldBindQuery(pg); err != nil {
		h.logger.Default().Error("pagination query binding err", zap.String("err", err.Error()))
		helpers.RespondError(c, err)
		return
	}

	res, err := h.uc.GenerateFeeds(c.Request.Context(), pg)

	if err != nil {
		//h.logger.Default().Error("generate feeds", zap.String("err", err.Error()))
		helpers.RespondError(c, err)
		return
	}

	helpers.RespondOK(c, res)
}
