package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	app := fiber.New()
	port := os.Getenv("PORT")

	app.Get("/", indexHandler)

	app.Post("/", postHandler)

	app.Put("/update", putHandler)

	app.Delete("/delete", deleteHandler)

	if port == "" {
		port = "3000"
	}

	psqlConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlConnection)
	CheckError(err)

	defer db.Close()
	err = db.Ping()

	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

// func main() {
// 	// connection string
// 	psqlConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

// 	// open database
// 	db, err := sql.Open("postgres", psqlConnection)
// 	CheckError(err)

// 	// close database
// 	defer db.Close()

// 	// check db
// 	err = db.Ping()
// 	CheckError(err)

// 	fmt.Println("Connected!")
// }

func CheckError(err error) {
	println(err)
	if err != nil {
		panic(err)
	}
}

func indexHandler(c *fiber.Ctx) error {
	return c.SendString("Hello")
}

func postHandler(c *fiber.Ctx) error {
	return c.SendString("Hello")
}

func putHandler(c *fiber.Ctx) error {
	return c.SendString("Hello")
}

func deleteHandler(c *fiber.Ctx) error {
	return c.SendString("Hello")
}
