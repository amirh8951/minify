package usecase

import (
	"context"
	"time"

	"minify/internal/domain"
)

// ShortenResult holds the output of a successful URL shortening.
type ShortenResult struct {
	ShortCode string
	ShortURL  string
	ExpiresIn string
}

// ShortenUseCase orchestrates the creation of a shortened URL.
type ShortenUseCase struct {
	repo    URLRepository
	baseURL string
	urlTTL  time.Duration
}

// NewShortenUseCase creates a new ShortenUseCase.
func NewShortenUseCase(repo URLRepository, baseURL string, urlTTL time.Duration) *ShortenUseCase {
	return &ShortenUseCase{
		repo:    repo,
		baseURL: baseURL,
		urlTTL:  urlTTL,
	}
}

// Execute validates the URL, generates a short code, and persists it.
func (uc *ShortenUseCase) Execute(ctx context.Context, rawURL string) (*ShortenResult, error) {
	originalURL, err := domain.NewOriginalURL(rawURL)
	if err != nil {
		return nil, err // domain.ErrInvalidURL
	}

	code, err := domain.GenerateShortCode()
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, code, originalURL, uc.urlTTL); err != nil {
		return nil, err
	}

	return &ShortenResult{
		ShortCode: code.String(),
		ShortURL:  uc.baseURL + "/" + code.String(),
		ExpiresIn: uc.urlTTL.String(),
	}, nil
}
