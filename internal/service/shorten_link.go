package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/nguyenthienan91/bookmark-manager/internal/repository"
	goredis "github.com/redis/go-redis/v9"
)

// ErrLinkNotFound is returned when the short code does not exist or has expired.
var ErrLinkNotFound = errors.New("short link not found")

const (
	shortCodeLength    = 7
	shortCodeCharset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCodeKeyPrefix = "shorturl:"
	maxRetryCount      = 10
)

// ShortenLink is the service interface for URL shortening.
//
//go:generate mockery --name ShortenLink --filename=shorten_link.go
type ShortenLink interface {
	// Shorten generates a unique 7-character code for the given URL.
	// expSeconds is the TTL in seconds; 0 means no expiration.
	Shorten(ctx context.Context, url string, expSeconds int) (string, error)

	// Resolve returns the original URL for the given short code.
	Resolve(ctx context.Context, code string) (string, error)
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

		key := fmt.Sprintf("%s%s", shortCodeKeyPrefix, code)
		created, err := s.repo.SaveLink(ctx, key, url, exp)
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

// Resolve looks up the original URL for the given short code.
func (s *shortenLinkService) Resolve(ctx context.Context, code string) (string, error) {
	key := fmt.Sprintf("%s%s", shortCodeKeyPrefix, code)
	url, err := s.repo.GetLink(ctx, key)
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", ErrLinkNotFound
		}
		return "", err
	}
	return url, nil
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
