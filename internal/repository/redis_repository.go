package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type urlRepository struct {
	client *redis.Client
}

func NewURLRepository(client *redis.Client) URLRepository {
	return &urlRepository{client: client}
}

func (r *urlRepository) Save(ctx context.Context, shortCode, originalURL string, ttl time.Duration) error {
	return r.client.Set(ctx, shortCode, originalURL, ttl).Err()
}

func (r *urlRepository) Get(ctx context.Context, shortCode string) (string, error) {
	val, err := r.client.Get(ctx, shortCode).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}
