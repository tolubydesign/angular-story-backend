package controller

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	dynamodb "github.com/tolubydesign/angular-story-backend/app/controller/dynamodb-request"
	postgres "github.com/tolubydesign/angular-story-backend/app/controller/postgres-request"
)

// Setup REST API request endpoints
func SetupMethods(app *fiber.App, db *sql.DB) {
	SetupPostgreSQLMethods(app, db)
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

// Setup REST API endpoints that use the PostgreSQL Database
func SetupPostgreSQLMethods(app *fiber.App, db *sql.DB) {
	app.Get("/stories", func(ctx *fiber.Ctx) error {
		return postgres.GetAllStoriesRequest(ctx, db)
	})

	app.Get("/story", func(ctx *fiber.Ctx) error {
		return postgres.GetSingleStoryRequest(ctx, db)
	})

	app.Post("/story", func(ctx *fiber.Ctx) error {
		return postgres.InsertStoryRequest(ctx, db)
	})

	app.Delete("/story", func(ctx *fiber.Ctx) error {
		return postgres.DeleteStoryRequest(ctx, db)
	})

	app.Put("/story", func(ctx *fiber.Ctx) error {
		return postgres.UpdateStoryRequest(ctx, db)
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return postgres.CheckHealthRequest(c, db)
	})
}

/*
Setup REST API endpoints that use the Dynamodb Database.
Connects to dynamodb database. If a connection is not made, process will error out.

Returns error if connection to dynamodb is incorrect/insufficient OR method cannot be created.
*/
func SetupDynamoDBMethods(app *fiber.App) {
	// Get Dynamodb Client
	client, err := dynamodb.ConnectDynamoDB()
	if err != nil {
		panic(err)
	}

	app.Get("/dynamodb-list-tables", func(ctx *fiber.Ctx) error {
		return dynamodb.GetAllDynamoDBTables(ctx, client)
	})

	app.Post("/dynamo-populate-database", func(ctx *fiber.Ctx) error {
		return dynamodb.PopulateDynamoDatabase(ctx, client)
	})

	app.Post("/dynamo-add-story", func(ctx *fiber.Ctx) error {
		return dynamodb.AddStoryRequest(ctx, client)
	})

	app.Get("/dynamo-get-story", func(ctx *fiber.Ctx) error {
		return dynamodb.GetStoryByIdRequest(ctx, client)
	})

	app.Get("/dynamo-list-stories", func(ctx *fiber.Ctx) error {
		return dynamodb.ListAllStories(ctx, client)
	})

	app.Put("/dynamo-update-story", func(ctx *fiber.Ctx) error {
		return dynamodb.UpdateDynamodbStoryRequest(ctx, client)
	})

	app.Delete("/dynamo-remove-story", func(ctx *fiber.Ctx) error {
		return dynamodb.DeleteDynamodbStoryRequest(ctx, client)
	})

	// Users
	app.Get("/dynamodb-list-users", func(ctx *fiber.Ctx) error {
		return dynamodb.ListAllUsersRequest(ctx, client)
	})
}
