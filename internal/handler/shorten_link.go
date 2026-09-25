package handler

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthienan91/bookmark-manager/internal/logger"
	"github.com/nguyenthienan91/bookmark-manager/internal/model"
	"github.com/nguyenthienan91/bookmark-manager/internal/service"
)

var validCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9]{1,20}$`)

// ShortenLink defines the handler interface for the link-shortening endpoint.
type ShortenLink interface {
	ShortenLink(c *gin.Context)
	RedirectLink(c *gin.Context)
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
		logger.L().Error().
			Err(err).
			Str("url", req.URL).
			Str("path", c.FullPath()).
			Msg("failed to shorten link")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ShortenLinkResponse{
		Code:    code,
		Message: "Shorten URL generated successfully!",
	})
}

// RedirectLink godoc
// @Summary      Redirect short link
// @Description  Redirect to the original URL for the given short code
// @Tags         Links
// @Param        code path string true "Shorten code"
// @Success      302
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /v1/links/redirect/{code} [get]
func (h *shortenLinkHandler) RedirectLink(c *gin.Context) {
	code := c.Param("code")
	if code == "" || !validCodeRegex.MatchString(code) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid or missing code"})
		return
	}

	originalURL, err := h.shortenLinkSvc.Resolve(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrLinkNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "short link not found"})
			return
		}
		logger.L().Error().
			Err(err).
			Str("code", code).
			Str("path", c.FullPath()).
			Msg("failed to resolve short link")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}
