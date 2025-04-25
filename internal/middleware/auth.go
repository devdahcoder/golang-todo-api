package middleware

import (
	"strings"
	"github.com/devdahcoder/golang-todo-api/pkg/token"

	"github.com/gofiber/fiber/v3"
)

func Auth(tokenMaker token.Maker) fiber.Handler {
    return func(c fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Authorization header is required",
            })
        }
        
        // Check if the header has the Bearer prefix
        fields := strings.Fields(authHeader)
        if len(fields) != 2 || fields[0] != "Bearer" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid authorization format. Use Bearer {token}",
            })
        }
        
        // Extract the token
        accessToken := fields[1]
        
        // Verify the token
        payload, err := tokenMaker.VerifyToken(accessToken)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid or expired token",
            })
        }
        
        // Store user ID in context for later use
        c.Locals("userId", payload.UserID)
        
        return c.Next()
    }
}