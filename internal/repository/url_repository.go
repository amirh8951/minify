package repository

import (
	"context"
	"time"
)

type URLRepository interface {
	Save(ctx context.Context, shortCode, originalURL string, ttl time.Duration) error
	Get(ctx context.Context, shortCode string) (string, error)
}
