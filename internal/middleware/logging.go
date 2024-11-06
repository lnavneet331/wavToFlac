package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// Store the request ID in context
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		c.Locals("requestID", requestID)

		// Process request
		err := c.Next()

		// Calculate response time
		responseTime := time.Since(start)

		// Log the request details
		fmt.Printf(
			"[%s] %s - %s %s - %d - %v\n",
			requestID,
			c.IP(),
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			responseTime,
		)

		return err
	}
}

func RateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Simple rate limiting implementation
		// In production, you might want to use Redis or a similar solution
		ip := c.IP()
		if isRateLimited(ip) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		}
		return c.Next()
	}
}

func isRateLimited(_ string) bool {
	// Implement rate limiting logic here
	// This is a placeholder implementation
	return false
}