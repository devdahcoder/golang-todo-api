package server

import (
	"github.com/devdahcoder/golang-todo-api/internal/config"
	"github.com/devdahcoder/golang-todo-api/internal/domain/user"
	"github.com/devdahcoder/golang-todo-api/internal/handlers"
	"github.com/devdahcoder/golang-todo-api/internal/repository/postgres"
	"github.com/devdahcoder/golang-todo-api/pkg/token"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)


func New(cfg *config.Config) *fiber.App {

	userRepository := postgres.NewUserRepository(cfg.Db)

	tokenMaker, err := token.NewJWTMaker(cfg.JWTSecret)

	if err != nil {
		panic(err)
	}

	userService := user.NewService(userRepository, tokenMaker)

	userHandler := handlers.NewUserHandler(userService)

	app := fiber.New(fiber.Config{
		StrictRouting: true,
		ErrorHandler:  customErrorHandler,		
	})

	app.Use(logger.New())
    app.Use(recover.New())
    app.Use(cors.New())

	setupRoutes(app, userHandler, tokenMaker)

	return app
	
}

func customErrorHandler(c fiber.Ctx, err error) error {
    code := fiber.StatusInternalServerError
    
    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
    }
    
    return c.Status(code).JSON(fiber.Map{
        "error": err.Error(),
    })
}