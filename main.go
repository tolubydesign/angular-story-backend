package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	configuration "github.com/tolubydesign/angular-story-backend/app/config"
	"github.com/tolubydesign/angular-story-backend/app/controller"
	"github.com/tolubydesign/angular-story-backend/app/database"
	"github.com/tolubydesign/angular-story-backend/app/logging"
)

func main() {
	// Setup project configuration
	config, err := configuration.BuildConfiguration()
	if err != nil {
		logging.Error(err.Error())
		panic(err)
	}

	c := config.Configuration
	logging.Event("ENVIRONMENT :::", c.Environment)

	DevelopmentAPI()
	// if helpers.IsLambda() /* OR c.Environment == "production" */ {
	// 	cdk.Run()
	// } else {
	// 	DevelopmentAPI()
	// }
}

// TODO: description
func SetupApplication() *fiber.App {
	return fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).SendString(err.Error())
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			// Return from handler
			return nil
		},
	})
}

// Developing int development mode.
// This function runs a local api, for testing purposes.
func DevelopmentAPI() {
	configuration, err := configuration.GetConfiguration()
	if err != nil {
		panic(err)
	}

	environmentPort := configuration.Configuration.Port
	env := configuration.Configuration.Environment

	app := SetupApplication()
	controller.HandleCORS(app, env)
	controller.SetupMethods(app)

	// get client
	client, err := database.GetDynamoSingleton()
	if err != nil {
		panic(err)
	}

	// Add dummy data if in development mode
	if (env == "development") && (client != nil) {
		// setup database with dummy data
		err := database.AddDummyData(client)
		if err != nil {
			panic(err)
		}
	}

	if environmentPort == "" {
		environmentPort = "2100"
	}

	log.Fatalln(app.Listen(fmt.Sprintf(":%v", environmentPort)))
}
