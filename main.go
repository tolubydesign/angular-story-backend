package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tolubydesign/angular-story-backend/app/controller"
	"github.com/tolubydesign/angular-story-backend/app/utils"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

func main() {
	var envs map[string]string
	envs, envErr := godotenv.Read(".env")
	gottenEnv := os.Getenv("PORT")
	environment := os.Getenv("ENV")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	_, postgresErr := utils.ConnectToPostgreSQLDatabase()
	if postgresErr != nil {
		panic(postgresErr)
	}

	_, redisErr := utils.ConnectToRedisDatabase()
	if redisErr != nil {
		panic(redisErr)
	}

	postgresDatabase, getPostgresErr := utils.GetPostgreSQLDatabaseSingleton()
	if getPostgresErr != nil {
		panic(getPostgresErr)
	}

	environmentPort := envs["PORT"]
	fmt.Printf("Port  = %v \n", environmentPort)
	fmt.Printf("env port  = %v \n", gottenEnv)

	app := SetupApplication()
	controller.HandleCORS(app, environment)
	controller.SetupMethods(app, postgresDatabase)

	if environmentPort == "" {
		environmentPort = "2100"
	}

	if postgresErr = postgresDatabase.Ping(); postgresErr != nil {
		panic(postgresErr)
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
