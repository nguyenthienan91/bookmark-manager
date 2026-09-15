package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthienan91/bookmark-manager/internal/model"
	"github.com/nguyenthienan91/bookmark-manager/internal/service"
)

const passwordLength = 12

type GenPass interface {
	GeneratePassword(c *gin.Context)
}

type genPassHandler struct {
	genPassService service.GenPass
}

func NewGenPass(GenPassSvc service.GenPass) GenPass {
	return &genPassHandler{
		genPassService: GenPassSvc,
	}
}

// GeneratePassword godoc
// @Summary      Generate a random password
// @Description  Generate a secure random password with 12 characters
// @Tags         Password
// @Produce      json
// @Success      200  {object}  model.GeneratePasswordResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /generate-password [get]
func (g *genPassHandler) GeneratePassword(c *gin.Context) {
	pass, err := g.genPassService.GeneratePassword(passwordLength)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to generate password"})
		return
	}
	c.JSON(http.StatusOK, model.GeneratePasswordResponse{Password: pass})
}
