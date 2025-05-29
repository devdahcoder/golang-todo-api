package errors

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func NewError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func InternalServerError(c fiber.Ctx, message string, err error) error {
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"code":    http.StatusInternalServerError,
		"message": message,
		"error":   err.Error(),
	})
}

func BadRequestError(c fiber.Ctx, message string, err error, structure map[string]any, details map[string]any) error {
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
		"code":    http.StatusBadRequest,
		"message": message,
		"error":   err.Error(),
		"request_example": structure,
		"details": details,
	})
}

func UnAuthorizedError(c fiber.Ctx, message string, err error) error {
	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
		"code":    http.StatusUnauthorized,
		"message": message,
		"error":   err.Error(),
	})
}

func NotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
		Err:     nil,
	}
}

// func UnAuthorizedError(message string) *AppError {
// 	return &AppError{
// 		Code:    http.StatusUnauthorized,
// 		Message: message,
// 		Err:     nil,
// 	}
// }

func (e *AppError) Error() string {
    if e.Err != nil {
        return e.Err.Error()
    }
    return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Is(target error) bool {
	return e.Code == target.(*AppError).Code
}

func (e *AppError) ErrorCode() int {
	return e.Code
}

func (e *AppError) ErrorMessage() string {
	return e.Message
}

func (e *AppError) ErrorResponse() map[string]interface{} {
	return map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	}
}

func (e *AppError) ErrorResponseWithError() map[string]interface{} {
	return map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
		"error":   e.Err.Error(),
	}
}

func (e *AppError) ErrorResponseWithStatus() map[string]interface{} {
	return map[string]interface
{}{
		"code":    e.Code,
		"message": e.Message,
	}
}

func (e *AppError) ErrorResponseWithCode() map[string]interface{} {
	return map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	}
}

func (e *AppError) ErrorResponseWithMessage() map[string]interface{} {
	return map[string]interface{}{
		"message": e.Message,
	}
}

func (e *AppError) ErrorResponseWithErrorMessage() map[string]interface{} {
	return map[string]interface{}{
		"error":   e.Err.Error(),
		"message": e.Message,
	}
}

