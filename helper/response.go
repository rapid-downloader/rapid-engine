package helper

import "github.com/gofiber/fiber/v2"

func Error(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusBadRequest).
		JSON(fiber.Map{
			"status":  "error",
			"message": message,
		})
}

func Success(ctx *fiber.Ctx, body interface{}) error {
	return ctx.Status(fiber.StatusOK).
		JSON(fiber.Map{
			"status": "ok",
			"data":   body,
		})
}

func Created(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusCreated).
		JSON(fiber.Map{
			"status": "ok",
		})
}
