package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/tolubydesign/angular-story-backend/app/config"
	"github.com/tolubydesign/angular-story-backend/app/controller"
	"github.com/tolubydesign/angular-story-backend/app/database"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

func main() {
	// Setup project configuration
	config, err := config.BuildConfiguration()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	// Connect to Redis database
	_, redisErr := database.ConnectToRedisDatabase()
	if redisErr != nil {
		panic(redisErr)
	}

	environmentPort := config.Configuration.Port
	env := config.Configuration.Environment
	fmt.Printf("PORT  = %v \n", environmentPort)
	fmt.Printf("ENV  = %v \n", env)

	app := SetupApplication()
	controller.HandleCORS(app, env)
	controller.SetupMethods(app)

	if environmentPort == "" {
		environmentPort = "2100"
	}

	log.Fatalln(app.Listen(fmt.Sprintf(":%v", environmentPort)))
}

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
