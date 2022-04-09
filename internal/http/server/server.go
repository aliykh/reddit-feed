package server

import (
	"fmt"
	postsHttp "github.com/aliykh/reddit-feed/internal/posts/delivery/http"
	"github.com/aliykh/reddit-feed/internal/posts/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

	"github.com/aliykh/log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/aliykh/reddit-feed/internal/config"
	"github.com/aliykh/reddit-feed/internal/posts/usecase"
	"github.com/aliykh/reddit-feed/pkg/customErrors"
	"github.com/aliykh/reddit-feed/pkg/helpers"
)

type Server struct {
	cfg    *config.Config
	logger *log.Factory
	router *gin.Engine
	dbClient *mongo.Client
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func New(cfg *config.Config, logger *log.Factory, dbClient *mongo.Client) (*Server, error) {

	sv := &Server{
		cfg:    cfg,
		logger: logger,
		dbClient: dbClient,
	}

	router := gin.Default()

	// configuring go-validator
	sv.setupValidators()

	// register swagger handler
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// set gin router as our default http server router
	sv.router = router

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	// register all endpoints in our application
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

	postRepo := repository.New(s.logger, s.dbClient)
	postsUC := usecase.New(s.logger, postRepo)
	postsHandlers := postsHttp.New(s.logger, postsUC)

	v1 := s.router.Group("/api/v1")

	postsHttp.RegisterHandlers(v1, postsHandlers)

}
