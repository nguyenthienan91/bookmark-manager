package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthienan91/bookmark-manager/internal/handler"
	"github.com/nguyenthienan91/bookmark-manager/internal/service"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}


type engine struct {
	app *gin.Engine
	cfg *Config
}

func NewEngine(cfg *Config) Engine {
	app := &engine{
		app: gin.Default(),
		cfg: cfg,
	}
	app.initRoutes()	 
	return app
	}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTp to test the API without starting the server
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

func (e *engine) initRoutes() {

	// Register the password generation route
	genPassService := service.NewGenPass()
	genPassHandler := handler.NewGenPass(genPassService)

	// Register the health check route
	healthCheckService := service.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceID)
	healthCheckHandler := handler.NewHealthCheck(healthCheckService)

	e.app.GET("/generate-password", genPassHandler.GeneratePassword)
	e.app.GET("/health-check", healthCheckHandler.HealthCheck)
}

