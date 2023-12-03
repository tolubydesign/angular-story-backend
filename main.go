package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	configuration "github.com/tolubydesign/angular-story-backend/app/config"
	"github.com/tolubydesign/angular-story-backend/app/controller"
	"github.com/tolubydesign/angular-story-backend/app/logging"
	"github.com/tolubydesign/angular-story-backend/cdk"
)

func main() {
	// Setup project configuration
	config, err := configuration.BuildConfiguration()
	logging.Event("ENVIRONMENT %s", config.Configuration.Environment)

	if err != nil {
		logging.Error(err.Error())
		panic(err)
	}

	if config.Configuration.Environment == "development" {
		APIDevelopment()
	}

	if config.Configuration.Environment == "production" {
		cdk.RunCDK()
	}
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

// TODO: description
func APIDevelopment() {
	configuration, err := configuration.GetConfiguration()
	if err != nil {
		panic(err)
	}

	environmentPort := configuration.Configuration.Port
	env := configuration.Configuration.Environment

	app := SetupApplication()
	controller.HandleCORS(app, env)
	controller.SetupMethods(app)

	if environmentPort == "" {
		environmentPort = "2100"
	}

	log.Fatalln(app.Listen(fmt.Sprintf(":%v", environmentPort)))
}
