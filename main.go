package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tolubydesign/angular-story-backend/app/controller"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

// TODO: redis logging system.
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	var envs map[string]string
	envs, err := godotenv.Read(".env")
	gottenEnv := os.Getenv("PORT")
	environment := os.Getenv("ENV")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connection := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)

	// Connect to database
	db, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	environmentPort := envs["PORT"]
	fmt.Printf("Port  = %v \n", environmentPort)
	fmt.Printf("env port  = %v \n", gottenEnv)

	app := SetupApplication()
	controller.HandleCORS(app, environment)
	controller.SetupMethods(app, db)

	if environmentPort == "" {
		environmentPort = "2100"
	}

	if err = db.Ping(); err != nil {
		panic(err)
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
