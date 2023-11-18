package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	dynamodb "github.com/tolubydesign/angular-story-backend/app/controller/dynamodb-request"
)

// Setup REST API request endpoints
func SetupMethods(app *fiber.App) {
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

	// Data
	app.Get("/list-tables", func(ctx *fiber.Ctx) error {
		return dynamodb.GetAllDynamoDBTables(ctx, client)
	})

	app.Post("/populate-database", func(ctx *fiber.Ctx) error {
		return dynamodb.PopulateDynamoDatabase(ctx, client)
	})

	app.Post("/add-story", func(ctx *fiber.Ctx) error {
		return dynamodb.AddStoryRequest(ctx, client)
	})

	app.Get("/get-story", func(ctx *fiber.Ctx) error {
		return dynamodb.GetStoryByIdRequest(ctx, client)
	})

	app.Get("/list-stories", func(ctx *fiber.Ctx) error {
		return dynamodb.ListAllStoriesRequest(ctx, client)
	})

	app.Put("/update-story", func(ctx *fiber.Ctx) error {
		return dynamodb.UpdateDynamodbStoryRequest(ctx, client)
	})

	app.Delete("/remove-story", func(ctx *fiber.Ctx) error {
		return dynamodb.DeleteDynamodbStoryRequest(ctx, client)
	})

	// Users
	app.Get("/list-users", func(ctx *fiber.Ctx) error {
		return dynamodb.ListAllUsersRequest(ctx, client)
	})

	// Users: login
	app.Get("/login", func(c *fiber.Ctx) error {
		return dynamodb.UserLoginRequest(c, client)
	})

	app.Post("/sign-up", func(c *fiber.Ctx) error {
		return dynamodb.UserSignUpRequest(c, client)
	})
}
