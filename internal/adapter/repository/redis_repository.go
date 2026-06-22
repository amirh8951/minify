package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"myfirstproject/internal/domain"
	"myfirstproject/internal/usecase"
)

type urlRepository struct {
	client *redis.Client
}

// NewURLRepository creates a new Redis-backed URLRepository adapter.
func NewURLRepository(client *redis.Client) usecase.URLRepository {
	return &urlRepository{client: client}
}

func (r *urlRepository) Save(ctx context.Context, code domain.ShortCode, originalURL domain.OriginalURL, ttl time.Duration) error {
	return r.client.Set(ctx, code.String(), originalURL.String(), ttl).Err()
}

func (r *urlRepository) Get(ctx context.Context, code domain.ShortCode) (domain.OriginalURL, error) {
	val, err := r.client.Get(ctx, code.String()).Result()
	if err == redis.Nil || val == "" {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return domain.NewOriginalURL(val)
}
