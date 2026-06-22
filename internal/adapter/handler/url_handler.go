package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"myfirstproject/internal/domain"
	"myfirstproject/internal/usecase"
)

type urlHandler struct {
	shorten  *usecase.ShortenUseCase
	redirect *usecase.RedirectUseCase
}

// NewURLHandler creates a new HTTP handler adapter.
func NewURLHandler(shorten *usecase.ShortenUseCase, redirect *usecase.RedirectUseCase) *urlHandler {
	return &urlHandler{
		shorten:  shorten,
		redirect: redirect,
	}
}

// request / response DTOs — adapter-local, not domain shared

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
	ExpiresIn string `json:"expires_in"`
}

type errorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Create handles POST /api/v1/shorten
func (h *urlHandler) Create(c fiber.Ctx) error {
	var req shortenRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Success: false,
			Message: "invalid request body",
		})
	}

	resp, err := h.shorten.Execute(c.Context(), req.URL)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidURL) {
			return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
				Success: false,
				Message: "invalid url",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(shortenResponse{
		ShortCode: resp.ShortCode,
		ShortURL:  resp.ShortURL,
		ExpiresIn: resp.ExpiresIn,
	})
}

// Redirect handles GET /:code
func (h *urlHandler) Redirect(c fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Success: false,
			Message: "missing short code",
		})
	}

	originalURL, err := h.redirect.Execute(c.Context(), code)
	if err != nil {
		if errors.Is(err, domain.ErrShortCodeNotFound) || errors.Is(err, domain.ErrEmptyShortCode) {
			return c.Status(fiber.StatusNotFound).JSON(errorResponse{
				Success: false,
				Message: "short code not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.Redirect().Status(fiber.StatusFound).To(originalURL)
}
