package server

import (
	"github.com/devdahcoder/golang-todo-api/internal/handlers"
	// "github.com/devdahcoder/golang-todo-api/internal/middleware"
	"github.com/devdahcoder/golang-todo-api/pkg/token"
	"github.com/gofiber/fiber/v3"
)

func setupRoutes(
    app *fiber.App,
    userHandler *handlers.UserHandler,
    tokenMaker token.Maker,
) {
    v1 := app.Group("/api/v1")
    
    auth := v1.Group("/auth")
    auth.Post("/register", userHandler.CreateUser)
    auth.Post("/login", userHandler.Login)
    
    users := v1.Group("/users")
    users.Get("/", userHandler.ListUsers)
    // users.Get("/:id", userHandler.GetUser)
    // users.Put("/:id", userHandler.UpdateUser)
    // users.Delete("/:id", userHandler.DeleteUser)
    
    app.Get("/health", func(c fiber.Ctx) error {
        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "status": "ok",
            "message": "API is running",
            "version": "1.0.0",
        })
    })

    
}

func customErrorHandler(c *fiber.Ctx, err error) error {
    // Default 500 statuscode
    // code := fiber.StatusInternalServerError
    
    // Check if it's a Fiber error
    // if e, ok := err.(*fiber.Error); ok {
    //     code = e.Code
    // }
    
    // Return JSON error response
    return nil
    // return c.Status(code).JSON(fiber.Map{
    //     "error": err.Error(),
    // })
}