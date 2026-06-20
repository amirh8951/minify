package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"

	"myfirstproject/internal/config"
	"myfirstproject/internal/handler"
	"myfirstproject/internal/middleware"
	"myfirstproject/internal/repository"
	"myfirstproject/internal/service"
	"myfirstproject/pkg/redis"
)

func main() {
	cfg := config.Load()

	rdb := redis.NewClient(cfg.RedisHost, cfg.RedisPort)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer rdb.Close()

	urlRepo := repository.NewURLRepository(rdb)
	urlSvc := service.NewURLService(urlRepo, cfg.BaseURL, cfg.URLTTL)
	urlHandler := handler.NewURLHandler(urlSvc)

	app := fiber.New()

	app.Use(middleware.Logger())

	app.Post("/api/v1/shorten", middleware.RateLimiter(rdb, cfg.RateLimit, cfg.RateLimitWindow), urlHandler.Create)
	app.Get("/:code", urlHandler.Redirect)

	go func() {
		addr := ":" + cfg.ServerPort
		log.Printf("server starting on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	if err := rdb.Close(); err != nil {
		log.Printf("redis close error: %v", err)
	}

	log.Println("server stopped")
}
