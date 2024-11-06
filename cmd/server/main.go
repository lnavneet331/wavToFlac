package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"audio-converter/internal/config"
	"audio-converter/internal/handlers"
	"audio-converter/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	// Middleware
	app.Use(middleware.Logger())

	// WebSocket route
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Routes
	app.Get("/health", handlers.HealthCheck)
	app.Get("/ws/convert", websocket.New(handlers.HandleAudioConversion))

	// Start server in a goroutine
	go func() {
		if err := app.Listen(cfg.ServerPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}