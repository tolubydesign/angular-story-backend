package utils

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/tolubydesign/angular-story-backend/app/models"
)

/*
Get request header id. From fiber context request.

Return string id
*/
func GetRequestHeaderIdWithoutValidation(c *fiber.Ctx) string {
	headers := c.GetReqHeaders()
	id := headers["Id"][0]
	return id
}

// Return id from the request header. Get request header id. From fiber context request.
// Errors will be noted if issues occur while trying to validate the provided id.
func GetRequestHeaderID(ctx *fiber.Ctx) (string, error) {
	var id string
	headers := ctx.GetReqHeaders()
	id = headers["Id"][0]

	err := ValidateLimitedStringVariable(id)
	if err != nil {
		return id, err
	}

	return id, err
}

var EndPoints = models.APIEndPoints{
	Get: models.GetEndpoint{
		Story: "/get-story",
		AllStories: "/list-stories",
		HealthCheck: "/health",
		Login: "/login",
		AllUsers: "/list-users",
	},
	Post: models.PostEndpoint{
		Story: "/add-story",
		PopulateDatabase: "/populate-database",
	},
	Put: models.PutEndpoint{
		Story: "/update-story",
	},
	Delete: models.DeleteEndpoint{
		Story: "/remove-story",
	},
}