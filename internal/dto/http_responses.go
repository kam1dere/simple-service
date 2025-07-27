package dto

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

// DTO  некоторых компаниях используется такой подход

const (
	FieldBadFormat     = "FIELD_BADFORMAT"
	FieldIncorrect     = "FIELD_INCORRECT"
	ServiceUnavailable = "SERVICE_UNAVAILABLE"
	InternalError      = "Service is currently unavailable. Please try again later."
	NotFound           = "NOT_FOUND"
)

type Response struct {
	Status string `json:"status"`
	Error  *Error `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

type Error struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

// TaskResponse — структура для ответа клиенту
type TaskResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NotFoundError(ctx *fiber.Ctx, desc string) error {
	return ctx.Status(fiber.StatusNotFound).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: NotFound,
			Desc: desc,
		},
	})
}

func BadResponseError(ctx *fiber.Ctx, code, desc string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: code,
			Desc: desc,
		},
	})
}

func InternalServerError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: ServiceUnavailable,
			Desc: InternalError,
		},
	})
}
