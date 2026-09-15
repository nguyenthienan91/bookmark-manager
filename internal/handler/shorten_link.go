package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthienan91/bookmark-manager/internal/model"
	"github.com/nguyenthienan91/bookmark-manager/internal/service"
)

// ShortenLink defines the handler interface for the link-shortening endpoint.
type ShortenLink interface {
	ShortenLink(c *gin.Context)
}

type shortenLinkHandler struct {
	shortenLinkSvc service.ShortenLink
}

// NewShortenLink creates a new ShortenLink handler.
func NewShortenLink(svc service.ShortenLink) ShortenLink {
	return &shortenLinkHandler{shortenLinkSvc: svc}
}

// ShortenLink godoc
// @Summary      Shorten a URL
// @Description  Generate a 7-character alphanumeric short code for the given URL
// @Tags         Links
// @Accept       json
// @Produce      json
// @Param        body body      model.ShortenLinkRequest  true  "Shorten Link Request"
// @Success      200  {object}  model.ShortenLinkResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /v1/links/shorten [post]
func (h *shortenLinkHandler) ShortenLink(c *gin.Context) {
	var req model.ShortenLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	if req.Exp < 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "exp must be non-negative"})
		return
	}

	code, err := h.shortenLinkSvc.Shorten(c.Request.Context(), req.URL, req.Exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ShortenLinkResponse{
		Code:    code,
		Message: "Shorten URL generated successfully!",
	})
}
