package model

// GeneratePasswordResponse represents the response for password generation
type GeneratePasswordResponse struct {
	Password string `json:"password"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}
