package controller

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupMethods(app *fiber.App, db *sql.DB) {
	app.Get("/stories", func(ctx *fiber.Ctx) error {
		return AllStoriesHandlerRequest(ctx, db)
	})

	app.Get("/story", func(ctx *fiber.Ctx) error {
		return RequestSingleStoryHandler(ctx, db)
	})

	app.Post("/story", func(ctx *fiber.Ctx) error {
		return InsertStory(ctx, db)
	})

	app.Delete("/story", func(ctx *fiber.Ctx) error {
		return DeleteStoryRequest(ctx, db)
	})

	app.Put("/story", func(ctx *fiber.Ctx) error {
		return UpdateStoryRequest(ctx, db)
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return CheckHealth(c, db)
	})
}

func HandleCORS(app *fiber.App, environment string) {
	// Initialize default config
	app.Use(cors.New())

	var configuration cors.Config
	if environment == "development" {
		configuration = cors.Config{
			AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
			AllowOrigins:     "*",
			AllowCredentials: true,
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
			// AllowOrigins: "https://gofiber.io, https://gofiber.net",
			// AllowHeaders: "Origin, Content-Type, Accept",
		}
	}

	// Or extend your config for customization
	app.Use(cors.New(configuration))
}
