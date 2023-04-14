package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

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

type Story struct {
	story_id    string
	title       string
	description string
	content     interface{}
}

type Service struct {
	db *sql.DB
}

func main() {
	var envs map[string]string
	envs, err := godotenv.Read(".env")
	gottenEnv := os.Getenv("PORT")
	// port := os.Getenv("PORT")
	// s3Bucket := os.Getenv("S3_BUCKET")
	// secretKey := os.Getenv("SECRET_KEY")
	// environment := envs["ENV"]
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connection := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)
	// connStr := "postgresql://<username>:<password>@<database_ip>/todos?sslmode=disable"
	// postgresConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// open database
	// db, err := sql.Open("postgres", postgresConnection)

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

// const (
// 	host     = "localhost"
// 	port     = 5400
// 	user     = "postgres"
// 	password = "man1234"
// 	dbname   = "DB_1"
// )

// func main() {
// 	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

// 	db, err := sql.Open("postgres", psqlconn)
// 	CheckError(err)

// 	defer db.Close()

// 	// insert
// 	// hardcoded
// 	insertStmt := `insert into "Students"("Name", "Roll_Number") values('Jacob', 20)`
// 	_, e := db.Exec(insertStmt)
// 	CheckError(e)

// 	// dynamic
// 	insertDynStmt := `insert into "Students"("Name", "Roll_Number") values($1, $2)`
// 	_, e = db.Exec(insertDynStmt, "Jack", 21)
// 	CheckError(e)
// }

// func CheckError(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

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
	requestDynamicString := `select * from Story`
	resp, err := fetchAllStories(requestDynamicString, db)

	if err != nil {
		ErrorResponse(err)
	}

	// return ctx.BodyParser(response)
	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	// ctx.Response().SetBody(fmt.Append(nil, resp))
	// ctx.Response().SetBody([]byte(fmt.Sprintf("%v", resp)))
	return ctx.SendString(fmt.Sprintf("%v", resp))
	// return ctx.SendString(string(response))

	// rows, err := db.Query(requestDynamicString)

	// if err != nil {
	// 	ErrorResponse(err)
	// 	// log.Fatal(err)
	// }

	// // An album slice to hold data from returned rows.
	// var stories []Story

	// // Loop through rows, using Scan to assign column data to struct fields.
	// for rows.Next() {
	// 	var story Story

	// 	if err := rows.Scan(&story.story_id, &story.description, &story.title, &story.content); err != nil {
	// 		return err
	// 	}

	// 	// fmt.Printf("rows - %d", story.content)
	// 	stories = append(stories, story)
	// }

	// fmt.Printf("rows: %v\n", rows)

	// // fmt.Printf(rows, "this is the second part of body\n")
	// if err = rows.Err(); err != nil {
	// 	return err
	// }

	// defer rows.Close()
	// return ctx.JSON(stories)

	// // then write more body
	// fmt.Fprintf(ctx, "this is the second part of body\n")

	// ctx.Response()
	// ctx.SendStatus(200)
	// ctx.Response().SetBody([]byte("rows"))
	// // return ctx.SendString(rows)

	// fmt.Printf("rows - ", rows)

	// return ctx.SendString("Hello Delete Handler")
}

func ErrorResponse(err error) {
	if err != nil {
		panic(err)
	}
}

func fetchAllStories(request string, db *sql.DB) ([]Story, error) {
	rows, err := db.Query(request)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var stories []Story

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var story Story
		if err := rows.Scan(&story.story_id, &story.title, &story.description, &story.content); err != nil {
			return stories, err
		}

		jsonData, err := json.Marshal(story)

		stories = append(stories, story)
		fmt.Printf("story: %v\n\n", story)
	}

	if err = rows.Err(); err != nil {
		return stories, err
	}

	return stories, nil
}
