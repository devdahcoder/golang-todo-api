package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	requestIDKey = "X-Request-ID"
)

func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get request ID from header or generate new one
		requestID := c.Get(requestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in header
		c.Set(requestIDKey, requestID)
		
		// Store in context for logging
		c.Locals("requestID", requestID)
		
		return c.Next()
	}
} 