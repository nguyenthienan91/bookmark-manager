package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthienan91/bookmark-manager/internal/model"
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
// @Description  Check if the service is running and Redis is reachable
// @Tags         Health
// @Produce      json
// @Success      200  {object}  model.HealthCheckResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /health-check [get]
func (h *healthCheckHandler) HealthCheck(c *gin.Context) {
	response, err := h.healthCheckService.Check(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
