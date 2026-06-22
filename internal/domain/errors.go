package domain

import "errors"

var (
	ErrInvalidURL       = errors.New("invalid url")
	ErrShortCodeNotFound = errors.New("short code not found")
	ErrEmptyShortCode   = errors.New("short code is empty")
)
