package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tolubydesign/angular-story-backend/app/models"
)

// Default error handler
var DefaultErrorHandler = func(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	// Return status code with error message
	return c.Status(code).SendString(err.Error())
}

func BasicErrorHandling(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	errorJson := &models.ErrorResponse{
		ErrorMessage: err.Error(),
		Code:         code,
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(code).JSON(errorJson)
}

// Override default error handler
var ErrorHandler = func(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Send custom error page
	err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
	if err != nil {
		// In case the SendFile fails
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	// Return from handler
	return nil
}

func RespondError(ctx *fiber.Ctx, responseCode int, errMessage interface{}) error {
	errorJson := &models.HTTPError{
		ErrorMessage: errMessage,
	}
	return ctx.Status(responseCode).JSON(errorJson)
}

func RespondSuccess(ctx *fiber.Ctx, responseCode int, data interface{}) error {
	return ctx.Status(responseCode).JSON(data)
}

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

type ResponseHandlerParameters struct {
	Context *fiber.Ctx
	Error   bool
	Code    int
	Message string
	Data    interface{}
}

// Basic method to handle request responses
func HandleResponse(context ResponseHandlerParameters) error {
	c := context.Context
	if c == nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint("Request context was not provided. Message: ", context.Message))
	}

	c.Response().StatusCode()
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(context.Code).JSON(models.HTTPResponse{
		Error:   context.Error,
		Code:    context.Code,
		Message: context.Message,
		Data:    context.Data,
	})
}
