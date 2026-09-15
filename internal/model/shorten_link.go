package model

// ShortenLinkRequest represents the request body for URL shortening.
type ShortenLinkRequest struct {
	URL string `json:"url" binding:"required,url"`
	Exp int    `json:"exp"`
}

// ShortenLinkResponse represents the response for URL shortening.
type ShortenLinkResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}