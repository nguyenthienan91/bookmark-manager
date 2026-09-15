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

// NewEngine creates a new Gin engine.
// shortenLinkSvc is optional; pass nil to disable the /v1/links/shorten route (e.g. in existing integration tests).
func NewEngine(cfg *Config, shortenLinkSvc ...service.ShortenLink) Engine {
	var svc service.ShortenLink
	if len(shortenLinkSvc) > 0 {
		svc = shortenLinkSvc[0]
	}
	app := &engine{
		app: gin.Default(),
		cfg: cfg,
	}
	app.initRoutes(svc)
	return app
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP lets tests call the router without starting a real server.
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

func (e *engine) initRoutes(shortenLinkSvc service.ShortenLink) {
	// Password generation
	genPassService := service.NewGenPass()
	genPassHandler := handler.NewGenPass(genPassService)

	// Health check
	healthCheckService := service.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceID)
	healthCheckHandler := handler.NewHealthCheck(healthCheckService)

	e.app.GET("/generate-password", genPassHandler.GeneratePassword)
	e.app.GET("/health-check", healthCheckHandler.HealthCheck)

	// Shorten link (only registered when a real service is wired in)
	if shortenLinkSvc != nil {
		shortenLinkHandler := handler.NewShortenLink(shortenLinkSvc)
		v1 := e.app.Group("/v1")
		v1.POST("/links/shorten", shortenLinkHandler.ShortenLink)
	}

	// Swagger docs
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

