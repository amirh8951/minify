package domain

import (
	"crypto/rand"
	"math/big"
	"net/url"
	"regexp"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var validShortCode = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// ShortCode is a value object representing a unique short code.
type ShortCode string

// NewShortCode creates a ShortCode value object. Returns ErrEmptyShortCode
// if the code is empty or contains invalid characters.
func NewShortCode(code string) (ShortCode, error) {
	if code == "" || !validShortCode.MatchString(code) {
		return "", ErrEmptyShortCode
	}
	return ShortCode(code), nil
}

// String returns the string representation.
func (s ShortCode) String() string { return string(s) }

// OriginalURL is a value object representing a validated URL.
type OriginalURL string

// NewOriginalURL parses and validates a URL string. Returns ErrInvalidURL
// if the string is not a valid HTTP/HTTPS URL.
func NewOriginalURL(rawURL string) (OriginalURL, error) {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", ErrInvalidURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", ErrInvalidURL
	}
	return OriginalURL(rawURL), nil
}

// String returns the string representation.
func (u OriginalURL) String() string { return string(u) }

// ShortenedURL is a domain entity representing a shortened URL mapping.
type ShortenedURL struct {
	ShortCode   ShortCode
	OriginalURL OriginalURL
	ExpiresAt   time.Time
}

// GenerateShortCode creates a cryptographically random 7-character alphanumeric code.
func GenerateShortCode() (ShortCode, error) {
	code := make([]byte, 7)
	for i := range code {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[idx.Int64()]
	}
	return NewShortCode(string(code))
}
