package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthienan91/bookmark-manager/internal/handler"
	"github.com/nguyenthienan91/bookmark-manager/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/nguyenthienan91/bookmark-manager/docs"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *Config
}

// NewEngine creates a new Gin engine with all application routes registered.
func NewEngine(cfg *Config, healthCheckSvc service.HealthCheck, shortenLinkSvc service.ShortenLink) Engine {
	app := &engine{
		app: gin.Default(),
		cfg: cfg,
	}
	app.initRoutes(healthCheckSvc, shortenLinkSvc)
	return app
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP lets tests call the router without starting a real server.
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

func (e *engine) initRoutes(healthCheckSvc service.HealthCheck, shortenLinkSvc service.ShortenLink) {
	// Password generation
	genPassService := service.NewGenPass()
	genPassHandler := handler.NewGenPass(genPassService)

	// Health check
	healthCheckHandler := handler.NewHealthCheck(healthCheckSvc)

	e.app.GET("/generate-password", genPassHandler.GeneratePassword)
	e.app.GET("/health-check", healthCheckHandler.HealthCheck)

	// Shorten link
	shortenLinkHandler := handler.NewShortenLink(shortenLinkSvc)
	v1 := e.app.Group("/v1")
	v1.POST("/links/shorten", shortenLinkHandler.ShortenLink)

	// Swagger docs
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
