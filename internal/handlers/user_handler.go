package handlers

import (
	er "errors"
	"strconv"

	"github.com/devdahcoder/golang-todo-api/internal/domain/user"
	"github.com/devdahcoder/golang-todo-api/pkg/errors"
	"github.com/devdahcoder/golang-todo-api/pkg/logger"
	"github.com/devdahcoder/golang-todo-api/pkg/validator"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
    userService user.Service
    ZapLogger *logger.Logger
}

func NewUserHandler(userService user.Service, zapLogger *logger.Logger) *UserHandler {
    return &UserHandler{
        userService: userService,
        ZapLogger: zapLogger,
    }
}

var (
	exampleEmail = "user@example.com"
	examplePassword = "password123"
	exampleUsername = "johndoe"
	signupRequestExample = fiber.Map{
		"email":    exampleEmail,
		"password": examplePassword,
		"username": exampleUsername,
		"first_name": "John",
		"last_name": "Doe",
	}
	signinRequestExample = fiber.Map{
		"email":    exampleEmail,
		"password": examplePassword,
	}
)

type queryParams struct {
	value map[string]string
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {
    id, err := strconv.ParseUint(c.Params("id"), 10, 32)

    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid user ID",
        })
    }
    
    ctx := c.Context()
    u, err := h.userService.GetUser(ctx, uint(id))
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
    
    return c.JSON(u)
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
    var input user.CreateUserInput
    
    // if err := c.BodyParser(&input); err != nil {
    //     return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    //         "error": "Invalid request body",
    //     })
    // }
    
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

func (h *UserHandler) Login(c fiber.Ctx) error {
    var input user.LoginInput

    if err := validator.InvalidFieldValidation(c, map[string]bool{
		"email":    true,
		"password": true,
	}, input); err != nil {
		if invalidFieldErr, ok := validator.IsInvalidFieldError(err); ok {
            return errors.BadRequestError(c, "Invalid request body", err, signinRequestExample, fiber.Map{
				"invalid_fields": invalidFieldErr.Fields,
			})
		}
		return errors.BadRequestError(c, "invalid request body", err, signinRequestExample, fiber.Map{})
	}

	v := validator.NewErrorValidator()
	v.Check(input.Email != "", "email", "email must be provided")
	v.Check(input.Password != "", "password", "password must be provided")

	if !v.IsValid() {
		return errors.BadRequestError(c, "invalid request body", er.New("invalid request body"), signinRequestExample, fiber.Map{
            "invalid_fields": v.ValidationErrorField,
        })
	}
    
    if err := c.Bind().Body(&input); err != nil {
        return errors.BadRequestError(c, "invalid request body", er.New("invalid request body"), signinRequestExample, fiber.Map{
            "invalid_fields": "invalid request body",
        })
    }
    
    ctx := c.Context()
    response, err := h.userService.Login(ctx, input)
    if err != nil {
        if err == user.ErrInvalidCredentials {
            return errors.UnAuthorizedError(c, "invalid credentials", err)
        }
        if err == user.ErrUserNotFound {
            return errors.UnAuthorizedError(c, "user not found", err)
        }
        return errors.InternalServerError(c, "login failed", err)
    }
    
    return c.JSON(response)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
    return nil
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
    return nil
}

func (h *UserHandler) ListUsers(c fiber.Ctx) error {
    ctx := c.Context()
    validator := validator.NewQueryValidator()
    
    rules := map[string]string{
        "limit":  "number",
        "offset":    "number",
    }
    
    if errors := validator.ValidateQuery(c, rules); len(errors) > 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "errors": errors,
        })
    }

    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    offset, _ := strconv.Atoi(c.Query("offset", "0"))
    if limit <= 0 || offset < 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid limit or offset",
        })
    }

    users, err := h.userService.ListUsers(ctx, limit, offset)

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to list users",
        })
    }

    return c.JSON(users)
}
