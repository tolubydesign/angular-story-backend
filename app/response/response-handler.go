package response

import "github.com/gofiber/fiber/v2"

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Story `json:"data"`
	Message string  `json:"message"`
}

// HTTPError struct to HTTPError object.
type HTTPError struct {
	ErrorMessage interface{} `json:"errorMessage"`
}

func RespondError(ctx *fiber.Ctx, responseCode int, errMessage interface{}) error {
	errorJson := &HTTPError{
		ErrorMessage: errMessage,
	}
	return ctx.Status(responseCode).JSON(errorJson)
}

func RespondSuccess(ctx *fiber.Ctx, responseCode int, data interface{}) error {
	return ctx.Status(responseCode).JSON(data)
}

func HandleDataResponse() {
	return
}
