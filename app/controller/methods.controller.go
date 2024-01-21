package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/tolubydesign/angular-story-backend/app/controller/dynamodb-request"
	"github.com/tolubydesign/angular-story-backend/app/config"
)

var client *dynamodb.Client

// Setup REST API request endpoints
func SetupMethods(app *fiber.App) {
	// Get Dynamodb Client
	client, err := dynamodbrequest.ConnectDynamoDB()
	if err != nil {
		message := fmt.Sprintf("Failed to connect with dynamodb database: %s", err.Error())
		panic(message)
	}

	SetupDynamoDBMethods(app, client)
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
func SetupDynamoDBMethods(app *fiber.App, client *dynamodb.Client) {
	configuration, err := config.GetConfiguration()
	if err != nil {
		message := fmt.Sprintf("Configuration failure: %s", err.Error())
		panic(message)
	}

	// Data
	app.Get("/list-tables", func(ctx *fiber.Ctx) error {
		return dynamodbrequest.GetAllDynamoDBTables(ctx, client, configuration)
	})

	// TODO: [decision] should this be allowed, as a method
	app.Post("/populate-database", func(ctx *fiber.Ctx) error {
		return dynamodbrequest.PopulateDynamoDatabase(ctx, client, configuration)
	})

	app.Post("/add-story", func(ctx *fiber.Ctx) error {
		return dynamodbrequest.AddStoryRequest(ctx, client, configuration)
	})

	app.Get("/get-story", func(ctx *fiber.Ctx) error {
		return dynamodbrequest.GetStoryByIdRequest(ctx, client, configuration)
	})

	app.Get("/list-stories", func(ctx *fiber.Ctx) error {
		return dynamodbrequest.ListAllStoriesRequest(ctx, client, configuration)
	})

	app.Put("/update-story", func(ctx *fiber.Ctx) error {
		return dynamodbrequest.UpdateDynamodbStoryRequest(ctx, client, configuration)
	})

	app.Delete("/remove-story", func(ctx *fiber.Ctx) error {
		return dynamodbrequest.DeleteDynamodbStoryRequest(ctx, client, configuration)
	})

	// Users
	app.Get("/list-users", func(ctx *fiber.Ctx) error {
		return dynamodbrequest.ListAllUsersRequest(ctx, client, configuration)
	})

	// Users: login
	app.Get("/login", func(c *fiber.Ctx) error {
		return dynamodbrequest.UserLoginRequest(c, client)
	})

	app.Post("/sign-up", func(c *fiber.Ctx) error {
		return dynamodbrequest.UserSignUpRequest(c, client, configuration)
	})

	// Health check
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return dynamodbrequest.HealthCheck(ctx, client, configuration)
	})
}
