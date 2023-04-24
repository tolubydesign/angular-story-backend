package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"

	"github.com/joho/godotenv"

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

type Service struct {
	db *sql.DB
}

type AllStories struct {
	Story []Story `json:"story"`
}

type Story struct {
	StoryId     uuid.UUID   `json:"story_id" validate:"uuid"`
	Title       string      `json:"title" validate:"required"`
	Description string      `json:"description" validate:"required"`
	Content     interface{} `json:"content"` // also of type StoryContent
}

type StoryContent struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Children    interface{} `json:"children"`
}

func main() {
	var envs map[string]string
	envs, err := godotenv.Read(".env")
	gottenEnv := os.Getenv("PORT")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connection := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)

	// Connect to database
	db, err := sql.Open("postgres", connection)

	environmentPort := envs["PORT"]
	fmt.Printf("Port  = %v \n", environmentPort)
	fmt.Printf("env port  = %v \n", gottenEnv)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Get("/all", func(c *fiber.Ctx) error {
		return requestAllHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})

	if environmentPort == "" {
		environmentPort = "2100"
	}

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("The database is connected")
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", environmentPort)))
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello Index Handler")
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello Post Handler")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello Put Handler")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello Delete Handler")
}

func requestAllHandler(ctx *fiber.Ctx, db *sql.DB) error {
	stories, err := fetchAllStories(db)

	var storyArray []Story
	for _, story := range stories {
		var content interface{}
		if story.Content != nil {
			str := fmt.Sprintf("%s", story.Content)
			byt := []byte(str)
			json.Unmarshal(byt, &content)
		}

		arrBookForPublic := Story{
			StoryId:     story.StoryId,
			Title:       story.Title,
			Description: story.Description,
			Content:     content,
		}
		storyArray = append(storyArray, arrBookForPublic)
	}

	allStories := AllStories{
		Story: storyArray,
	}

	if err != nil {
		panic(err)
	}

	// return ctx.BodyParser(response)
	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(allStories)
}

func ErrorResponse(err error) {
	if err != nil {
		panic(err)
	}
}

func fetchAllStories(db *sql.DB) ([]Story, error) {
	request := `select * from Story`
	rows, err := db.Query(request)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stories []Story

	for rows.Next() {
		var story Story
		err := rows.Scan(&story.StoryId, &story.Title, &story.Description, &story.Content)
		if err != nil {
			return stories, err
		}

		stories = append(stories, story)
	}

	if err = rows.Err(); err != nil {
		return stories, err
	}

	return stories, nil
}
