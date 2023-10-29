package controller

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupMethods(app *fiber.App, db *sql.DB) {
	app.Get("/stories", func(ctx *fiber.Ctx) error {
		return GetAllStoriesRequest(ctx, db)
	})

	app.Get("/story", func(ctx *fiber.Ctx) error {
		return GetSingleStoryRequest(ctx, db)
	})

	app.Post("/story", func(ctx *fiber.Ctx) error {
		return InsertStoryRequest(ctx, db)
	})

	app.Delete("/story", func(ctx *fiber.Ctx) error {
		return DeleteStoryRequest(ctx, db)
	})

	app.Put("/story", func(ctx *fiber.Ctx) error {
		return UpdateStoryRequest(ctx, db)
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return CheckHealthRequest(c, db)
	})

	SetupDynamoDBMethods(app)
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

/*
Setup dynamodb relate methods and requests.
Connects to dynamodb database. If a connection is not made, process will error out.

@returns - Error if connection to dynamodb is incorrect/insufficient OR method cannot be created.
*/
func SetupDynamoDBMethods(app *fiber.App) {
	// Get client
	client, err := ConnectDynamoDB()
	if err != nil {
		panic(err)
	}

	// app.Post("/dynamo-table", func(ctx *fiber.Ctx) error {
	// 	return CreateNewTable(ctx, client)
	// })

	app.Get("/dynamodb-list-tables", func(ctx *fiber.Ctx) error {
		return GetAllDynamoDBTables(ctx, client)
	})

	app.Post("/dynamo-populate-database", func(ctx *fiber.Ctx) error {
		return PopulateDynamoDatabase(ctx, client)
	})

	app.Post("/dynamo-add-story", func (ctx *fiber.Ctx) error {
		return AddStoryRequest(ctx, client)
	})

	app.Get("/dynamo-get-story", func (ctx *fiber.Ctx) error {
		return GetStoryByIdRequest(ctx, client)
	})

	app.Get("/dynamo-list-stories", func(ctx *fiber.Ctx) error {
		return ListAllStories(ctx, client)
	})
}
