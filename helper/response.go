package helper

import "github.com/gofiber/fiber/v2"

func Error(ctx *fiber.Ctx, message string, code ...int) error {
	statusCode := fiber.StatusBadRequest
	if len(code) > 0 {
		statusCode = code[0]
	}

	return ctx.Status(statusCode).
		JSON(fiber.Map{
			"message": message,
		})
}

func Success(ctx *fiber.Ctx, body ...interface{}) error {
	if len(body) == 0 {
		return ctx.SendStatus(200)
	}

	return ctx.Status(fiber.StatusOK).JSON(body[0])
}

func Created(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusCreated)
}

func NotFound(ctx *fiber.Ctx) error {
	return Error(ctx, "Not Found", fiber.StatusNotFound)
}

func Unauthorized(ctx *fiber.Ctx) error {
	return Error(ctx, "Unauthorized", fiber.StatusUnauthorized)
}
