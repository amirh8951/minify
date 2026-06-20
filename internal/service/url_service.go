package service

import (
	"context"

	"myfirstproject/internal/model"
)

type URLService interface {
	CreateShortURL(ctx context.Context, originalURL string) (*model.ShortenResponse, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
}
