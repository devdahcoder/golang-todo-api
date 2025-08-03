package server

import (
	"time"

	"github.com/devdahcoder/golang-todo-api/internal/config"
	"github.com/devdahcoder/golang-todo-api/internal/domain/user"
	"github.com/devdahcoder/golang-todo-api/internal/handlers"
	"github.com/devdahcoder/golang-todo-api/internal/middleware"
	"github.com/devdahcoder/golang-todo-api/internal/repository/postgres"
	"github.com/devdahcoder/golang-todo-api/pkg/token"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	// "github.com/gofiber/fiber/v3/middleware/timeout"
	"go.uber.org/zap"
)

func NewServer(cfg *config.Config) *fiber.App {
	userRepository := postgres.NewUserRepository(cfg.Db)

	tokenMaker, err := token.NewJWTMaker(cfg.JWTSecret)
	if err != nil {
		cfg.ZapLogger.Fatal("Failed to create token maker", zap.Error(err))
		panic(err)
	}

	userService := user.NewService(userRepository, tokenMaker)
	userHandler := handlers.NewUserHandler(userService, cfg.ZapLogger)

	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	app := fiber.New(fiber.Config{
		StrictRouting: true,
	})

	// Global middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(middleware.RequestID())
	app.Use(rateLimiter.Middleware())
	app.Use(rateLimiter.Middleware())
    // app.Use(timeout.New(timeout.Config{
    //     Timeout: 10 * time.Second,
    // }))

	// Setup routes
	setupRoutes(app, userHandler, tokenMaker)

	return app
}

// func customErrorHandler(c fiber.Ctx, err error) error {
//     code := fiber.StatusInternalServerError
    
//     if e, ok := err.(*fiber.Error); ok {
//         code = e.Code
//     }
    
//     return c.Status(code).JSON(fiber.Map{
//         "error": err.Error(),
//     })
// }