package usecase

import (
	"context"
	"time"

	"myfirstproject/internal/domain"
)

// URLRepository is the output port for persisting shortened URLs.
type URLRepository interface {
	Save(ctx context.Context, code domain.ShortCode, originalURL domain.OriginalURL, ttl time.Duration) error
	Get(ctx context.Context, code domain.ShortCode) (domain.OriginalURL, error)
}
