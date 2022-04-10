package http

import (
	"github.com/aliykh/reddit-feed/internal/posts"
	"github.com/gin-gonic/gin"
)

const path = "/post"

func RegisterHandlers(router *gin.RouterGroup, handlers posts.Handlers)  {

	r1Group := router.Group(path)
	r1Group.POST("/", handlers.Create)

	r1Group.GET("/generate", handlers.Generate)

}
