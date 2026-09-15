package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthienan91/bookmark-manager/internal/service"
)

type HealthCheck interface {
	HealthCheck(c *gin.Context)
}

type healthCheckHandler struct {
	healthCheckService service.HealthCheck
}

func NewHealthCheck(healthCheckSvc service.HealthCheck) HealthCheck {
	return &healthCheckHandler{
		healthCheckService: healthCheckSvc,
	}
}

// HealthCheck godoc
// @Summary      Health check endpoint
// @Description  Check if the service is running and return service information
// @Tags         Health
// @Produce      json
// @Success      200  {object}  model.HealthCheckResponse
// @Router       /health-check [get]
func (h *healthCheckHandler) HealthCheck(c *gin.Context) {
	response := h.healthCheckService.Check()

	c.JSON(http.StatusOK, response)
}