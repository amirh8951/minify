package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"myfirstproject/internal/model"
	"myfirstproject/internal/service"
)

type URLHandler struct {
	svc service.URLService
}

func NewURLHandler(svc service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

func (h *URLHandler) Create(c fiber.Ctx) error {
	var req model.ShortenRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Message: "invalid request body",
		})
	}

	ctx := context.Background()
	resp, err := h.svc.CreateShortURL(ctx, req.URL)
	if err != nil {
		if err.Error() == "invalid url" {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Success: false,
				Message: "invalid url",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *URLHandler) Redirect(c fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Message: "missing short code",
		})
	}

	ctx := context.Background()
	originalURL, err := h.svc.GetOriginalURL(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Success: false,
			Message: "internal server error",
		})
	}

	if originalURL == "" {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Success: false,
			Message: "short code not found",
		})
	}

	return c.Redirect().Status(fiber.StatusFound).To(originalURL)
}
