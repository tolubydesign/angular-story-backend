package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/tolubydesign/angular-story-backend/app/models"
	"github.com/tolubydesign/angular-story-backend/app/queries"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

func RequestAllStoriesHandler(ctx *fiber.Ctx, db *sql.DB) error {
	stories, err := queries.GetAllStories(db)

	var storyArray []models.Story
	for _, story := range stories {
		var content interface{}
		if story.Content != nil {
			str := fmt.Sprintf("%s", story.Content)
			byt := []byte(str)
			json.Unmarshal(byt, &content)
		}

		returningStoryModel := models.Story{
			StoryId:     story.StoryId,
			Title:       story.Title,
			Description: story.Description,
			Content:     content,
		}
		storyArray = append(storyArray, returningStoryModel)
	}

	response := models.JSONResponse{
		Type:    "success",
		Data:    storyArray,
		Message: "",
	}

	if err != nil {
		panic(err)
	}

	// return ctx.BodyParser(response)
	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

func RequestSingleStoryHandler(ctx *fiber.Ctx, db *sql.DB) error {
	context := ctx.Context()
	headers := ctx.GetReqHeaders()
	headerId := headers["Id"]
	story, err := queries.GetSingleStory(headerId, context, db)

	if err != nil {
		panic(err)
	}

	var content interface{}
	if story.Content != nil {
		str := fmt.Sprintf("%s", story.Content)
		byt := []byte(str)
		json.Unmarshal(byt, &content)
	}

	returningResponse := models.Story{
		StoryId:     story.StoryId,
		Title:       story.Title,
		Description: story.Description,
		Content:     content,
	}

	response := models.JSONResponse{
		Type:    "success",
		Data:    returningResponse,
		Message: "",
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

func CheckHealth(c *fiber.Ctx, db *sql.DB) error {
	ctxt := c.Context()
	ctx, cancel := context.WithTimeout(ctxt, 2*time.Second)

	defer cancel()
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	c.Response().StatusCode()
	return c.SendString("Database is working.")
}

func InsertStory(c *fiber.Ctx, db *sql.DB) error {
	// TODO: verify that information being sent in is valid json
	fiberContext := c.Context()
	ctx, cancel := context.WithTimeout(fiberContext, 3*time.Second)
	defer cancel()

	var body models.Story
	byteBody := c.Body()
	json.Unmarshal(byteBody, &body)
	bodyContentJson, err := json.Marshal(body.Content)
	if err != nil {
		panic(err)
	}

	model := models.Story{
		Title:       body.Title,
		Description: body.Description,
		Content:     bodyContentJson,
	}

	execution := "INSERT INTO story(title, description, content) VALUES($1, $2, $3);"
	result, err := db.ExecContext(ctx, execution, model.Title, model.Description, model.Content)

	if err != nil {
		log.Fatalf("Fatal Results Error - %s", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Fatal Rows Error: %d", err)
	}
	if rows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", rows)
	}

	response := models.JSONResponse{
		Type:    "success",
		Data:    nil,
		Message: "Database has been updated.",
	}

	c.Response().StatusCode()
	c.Response().Header.Add("Content-Type", "application/json")
	return c.JSON(response)
}
