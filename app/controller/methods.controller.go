package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/tolubydesign/angular-story-backend/app/config"
	dynamodbrequest "github.com/tolubydesign/angular-story-backend/app/controller/dynamodb-request"
	"github.com/tolubydesign/angular-story-backend/app/utils"
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
	// "/list-tables"
	app.Get(utils.Endpoints.Get.Tables, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.GetAllDynamoDBTables(ctx, client, configuration)
	})

	// For local development
	// "/populate-database"
	// app.Post(utils.Endpoints.Post.PopulateDatabase, func(ctx *fiber.Ctx) error {
	// 	return dynamodbrequest.PopulateDynamoDatabase(ctx, client, configuration)
	// })

	// "/add-story"
	app.Post(utils.Endpoints.Post.Story, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.AddStoryRequest(ctx, client, configuration)
	})

	// "/get-story"
	app.Get(utils.Endpoints.Get.Story, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.GetStoryByIdRequest(ctx, client, configuration)
	})

	// "/list-stories"
	app.Get(utils.Endpoints.Get.AllStories, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.ListAllStoriesRequest(ctx, client, configuration)
	})

	// "/update-story"
	app.Put(utils.Endpoints.Put.Story, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.UpdateDynamodbStoryRequest(ctx, client, configuration)
	})

	// "/remove-story"
	app.Delete(utils.Endpoints.Delete.Story, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.DeleteDynamodbStoryRequest(ctx, client, configuration)
	})

	// Users
	// "/list-users"
	app.Get(utils.Endpoints.Get.Users, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.ListAllUsersRequest(ctx, client, configuration)
	})

	// Users: login
	// "/login"
	app.Get(utils.Endpoints.Get.Login, func(c *fiber.Ctx) error {
		return dynamodbrequest.UserLoginRequest(c, client)
	})

	// "/sign-up"
	app.Post(utils.Endpoints.Post.SignUp, func(c *fiber.Ctx) error {
		return dynamodbrequest.UserSignUpRequest(c, client, configuration)
	})

	// Health check
	// "/health"
	app.Get(utils.Endpoints.Get.HealthCheck, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.HealthCheck(ctx, client, configuration)
	})
}
