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

const (
	GetRequestEndpoint_Story       = "/story"
	GetRequestEndpoint_AllStories  = "/stories"
	GetRequestEndpoint_HealthCheck = "/health"
)

// A dictionary of Api endpoints available
var Endpoints = models.APIEndpoints{
	Get: models.GetEndpoint{
		Story:       "/get-story",
		AllStories:  "/list-stories",
		HealthCheck: "/health",
		Login:       "/login",
		Users:       "/list-users",
		Tables:      "/list-tables",
	},
	Post: models.PostEndpoint{
		Story:            "/add-story",
		PopulateDatabase: "/populate-database",
		SignUp:           "/sign-up",
	},
	Put: models.PutEndpoint{
		Story: "/update-story",
	},
	Delete: models.DeleteEndpoint{
		Story: "/remove-story",
	},
}
