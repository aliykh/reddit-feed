package server

import (

	"fmt"
	"github.com/aliykh/reddit-feed/pkg/customErrors"
	"github.com/aliykh/reddit-feed/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"net/http"

	"github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/config"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/go-playground/locales/en"
)

type Server struct {
	cfg    *config.Config
	logger *log.Factory
	router *gin.Engine
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func New(cfg *config.Config, logger *log.Factory) (*Server, error) {

	sv := &Server{
		cfg:    cfg,
		logger: logger,
	}

	router := gin.Default()

	sv.setupValidators()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	sv.router = router

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	sv.mapHandlers()

	return sv, nil
}

func (s *Server) setupValidators() {
	binding.Validator = new(helpers.DefaultValidator)

	engine := binding.Validator.Engine().(*validator.Validate)

	eng := en.New()
	uni := ut.New(eng, eng)
	customErrors.Trans, _ = uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(engine, customErrors.Trans)

}

func (s *Server) mapHandlers() {


	s.router.GET("/hello", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"result": "hey you!"})
	})

}
