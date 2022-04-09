package posts

import "github.com/gin-gonic/gin"

type Handlers interface {
	Create(c *gin.Context)
}