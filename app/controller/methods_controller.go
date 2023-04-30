package controller

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func SetupMethods(app *fiber.App, db *sql.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return IndexHandler(c, db)
	})

	app.Get("/all", func(ctx *fiber.Ctx) error {
		return RequestAllStoriesHandler(ctx, db)
	})

	app.Get("/stories", func(ctx *fiber.Ctx) error {
		return RequestSingleStoryHandler(ctx, db)
	})

	app.Post("/story", func(ctx *fiber.Ctx) error {
		return InsertStory(ctx, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return PostHandler(c, db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return PutHandler(c, db)
	}, func(c *fiber.Ctx) error {
		return c.SendString(c.Params("id"))
	})

	app.Delete("/delete", func(c *fiber.Ctx) error {
		return DeleteHandler(c, db)
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return CheckHealth(c, db)
	})

}
