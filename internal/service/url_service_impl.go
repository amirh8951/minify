package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"net/url"
	"time"

	"myfirstproject/internal/model"
	"myfirstproject/internal/repository"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type urlService struct {
	repo    repository.URLRepository
	baseURL string
	urlTTL  time.Duration
}

func NewURLService(repo repository.URLRepository, baseURL string, urlTTL time.Duration) URLService {
	return &urlService{
		repo:    repo,
		baseURL: baseURL,
		urlTTL:  urlTTL,
	}
}

func (s *urlService) CreateShortURL(ctx context.Context, originalURL string) (*model.ShortenResponse, error) {
	if !isValidURL(originalURL) {
		return nil, errors.New("invalid url")
	}

	code, err := generateShortCode()
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, code, originalURL, s.urlTTL); err != nil {
		return nil, err
	}

	return &model.ShortenResponse{
		ShortCode: code,
		ShortURL:  s.baseURL + "/" + code,
		ExpiresIn: s.urlTTL.String(),
	}, nil
}

func (s *urlService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	return s.repo.Get(ctx, shortCode)
}

func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func generateShortCode() (string, error) {
	code := make([]byte, 7)
	for i := range code {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[idx.Int64()]
	}
	return string(code), nil
}
