package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/nguyenthienan91/bookmark-manager/internal/repository"
)

const (
	shortCodeLength  = 7
	shortCodeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxRetryCount    = 10
)

// ShortenLink is the service interface for URL shortening.
//
//go:generate mockery --name ShortenLink --filename=shorten_link.go
type ShortenLink interface {
	// Shorten generates a unique 7-character code for the given URL.
	// expSeconds is the TTL in seconds; 0 means no expiration.
	Shorten(ctx context.Context, url string, expSeconds int) (string, error)
}

type shortenLinkService struct {
	repo repository.LinkRepository
}

// NewShortenLink creates a new ShortenLink service.
func NewShortenLink(repo repository.LinkRepository) ShortenLink {
	return &shortenLinkService{repo: repo}
}

func (s *shortenLinkService) Shorten(ctx context.Context, url string, expSeconds int) (string, error) {
	var exp time.Duration
	if expSeconds > 0 {
		exp = time.Duration(expSeconds) * time.Second
	}

	for i := 0; i < maxRetryCount; i++ {
		code, err := generateShortCode(shortCodeLength)
		if err != nil {
			return "", err
		}

		created, err := s.repo.SaveLink(ctx, code, url, exp)
		if err != nil {
			return "", err
		}
		if created {
			return code, nil
		}
		// code collision — retry with a new code
	}

	return "", fmt.Errorf("failed to generate a unique short code after %d retries", maxRetryCount)
}

// generateShortCode returns a cryptographically random alphanumeric string of the given length.
func generateShortCode(length int) (string, error) {
	code := make([]byte, length)
	charsetLen := big.NewInt(int64(len(shortCodeCharset)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		code[i] = shortCodeCharset[n.Int64()]
	}
	return string(code), nil
}
