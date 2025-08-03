package middleware

import (
	"strings"

	"github.com/devdahcoder/golang-todo-api/internal/domain/user"
	"github.com/devdahcoder/golang-todo-api/pkg/errors"
	"github.com/devdahcoder/golang-todo-api/pkg/token"
	"github.com/gofiber/fiber/v3"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(c fiber.Ctx) error {
		authorizationHeader := c.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			return errors.UnAuthorizedError(c, "authorization header is not provided", nil)
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			return errors.UnAuthorizedError(c, "invalid authorization header format", nil)
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return errors.UnAuthorizedError(c, "unsupported authorization type", nil)
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return errors.UnAuthorizedError(c, "invalid token", err)
		}

		// Store the payload in the context for later use
		c.Locals(authorizationPayloadKey, payload)
		return c.Next()
	}
}

// GetAuthPayload retrieves the authorization payload from the context
func GetAuthPayload(c fiber.Ctx) (*token.Payload, error) {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return payload, nil
}