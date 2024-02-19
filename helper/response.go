package helper

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/log"
)

func get(err ...error) error {
	var e error
	if len(err) > 0 {
		e = err[0]
	}

	return e
}

func Error(ctx *fiber.Ctx, code int, err ...error) error {
	e := get(err...)

	log.Println(e.Error())
	return ctx.Status(code).JSON(fiber.Map{
		"message": e.Error(),
	})
}

func Ok(ctx *fiber.Ctx, body ...interface{}) error {
	return Success(ctx, fiber.StatusOK, body...)
}

func Success(ctx *fiber.Ctx, code int, body ...interface{}) error {
	if len(body) == 0 {
		return ctx.SendStatus(code)
	}

	return ctx.Status(fiber.StatusOK).JSON(body[0])
}

func Created(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusCreated)
}

func NotFound(ctx *fiber.Ctx) error {
	return Error(ctx, fiber.StatusNotFound, fmt.Errorf("not found"))
}

func Unauthorized(ctx *fiber.Ctx) error {
	return Error(ctx, fiber.StatusUnauthorized, fmt.Errorf("unauthorized"))
}

func BadRequest(ctx *fiber.Ctx, err ...error) error {
	return Error(ctx, fiber.StatusBadRequest, err...)
}

func InternalServerError(ctx *fiber.Ctx, err ...error) error {
	return Error(ctx, fiber.StatusInternalServerError, err...)
}
