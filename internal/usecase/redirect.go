package usecase

import (
	"context"

	"minify/internal/domain"
)

// RedirectUseCase orchestrates resolving a short code to its original URL.
type RedirectUseCase struct {
	repo URLRepository
}

// NewRedirectUseCase creates a new RedirectUseCase.
func NewRedirectUseCase(repo URLRepository) *RedirectUseCase {
	return &RedirectUseCase{repo: repo}
}

// Execute looks up a short code and returns the original URL.
func (uc *RedirectUseCase) Execute(ctx context.Context, code string) (string, error) {
	shortCode, err := domain.NewShortCode(code)
	if err != nil {
		return "", err // domain.ErrEmptyShortCode
	}

	originalURL, err := uc.repo.Get(ctx, shortCode)
	if err != nil {
		return "", err
	}

	if originalURL == "" {
		return "", domain.ErrShortCodeNotFound
	}

	return originalURL.String(), nil
}
