package handlers

import (
	"strconv"

	"github.com/devdahcoder/golang-todo-api/internal/domain/user"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
    userService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid user ID",
        })
    }
    
    ctx := c.Context()
    user, err := h.userService.GetUser(ctx, uint(id))
    if err != nil {
        if err == user.ErrUserNotFound {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "User not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to get user",
        })
    }
    
    return c.JSON(user)
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
    var input user.CreateUserInput
    
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    // Validate input
    // Use a validation package here
    
    ctx := c.Context()
    createdUser, err := h.userService.CreateUser(ctx, input)
    if err != nil {
        if err == user.ErrEmailAlreadyExists {
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "error": "Email already exists",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create user",
        })
    }
    
    return c.Status(fiber.StatusCreated).JSON(createdUser)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
    var input user.LoginInput
    
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    ctx := c.Context()
    authResponse, err := h.userService.Login(ctx, input)
    if err != nil {
        if err == user.ErrInvalidCredentials {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid credentials",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Login failed",
        })
    }
    
    return c.JSON(authResponse)
}

// Implement other handler methods...