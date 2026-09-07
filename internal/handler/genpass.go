package handler

import (
	"github.com/gin-gonic/gin"
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

func (g *genPassHandler) GeneratePassword(c *gin.Context) {
	pass, err := g.genPassService.GeneratePassword(passwordLength)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate password"})
		return
	}
	c.JSON(200, gin.H{"password": pass})
}
