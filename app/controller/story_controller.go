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
	"github.com/tolubydesign/angular-story-backend/app/utils"
	"github.com/tolubydesign/angular-story-backend/pkg/response"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

func AllStoriesHandlerRequest(ctx *fiber.Ctx, db *sql.DB) error {
	stories, err := queries.GetAllStories(db)

	redisErr := utils.ConsoleActionToRedisDatabase("Attempting to get all stories")
	if redisErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, redisErr.Error())
	}

	var storyArray []models.Story
	for _, story := range stories {
		var content interface{}
		if story.Content != nil {
			// NOTE: look at https://go.dev/src/cmd/vet/testdata/print/print.go and https://pkg.go.dev/fmt#hdr-Printing
			// to address error "fmt.Sprintf format %s has arg story.Content of wrong type"
			// PREVIOUSLY: str := fmt.Sprintf("%s", story.Content)
			str := fmt.Sprintf("%s", story.Content)
			byt := []byte(str)
			json.Unmarshal(byt, &content)
		}

		returningStoryModel := models.Story{
			Id:          story.Id,
			Title:       story.Title,
			Description: story.Description,
			Content:     content,
		}
		storyArray = append(storyArray, returningStoryModel)
	}

	response := models.JSONResponse{
		Type:    "success",
		Data:    storyArray,
		Message: "Fetch all stories.",
	}

	if err != nil {
		panic(err)
	}

	redisErr = utils.ConsoleActionToRedisDatabase("Request to get all stories successful")
	if redisErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, redisErr.Error())
	}

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

	// Convert JSON to Struct
	var content interface{}
	if story.Content != nil {
		// NOTE: look at https://go.dev/src/cmd/vet/testdata/print/print.go and https://pkg.go.dev/fmt#hdr-Printing
		// to address error "fmt.Sprintf format %s has arg story.Content of wrong type"
		// PREVIOUSLY: str := fmt.Sprintf("%S", story.Content)
		str := []byte(fmt.Sprintf("%s", story.Content))
		byt := []byte(str)
		json.Unmarshal(byt, &content)
	}

	returningResponse := models.Story{
		Id:          story.Id,
		Title:       story.Title,
		Description: story.Description,
		Content:     content,
	}

	response := models.JSONResponse{
		Type:    "success",
		Data:    returningResponse,
		Message: "Fetch single story.",
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
	err := queries.AddStory(c, db)
	if err != nil {
		log.Fatalf("Insert Story Fatal Error - %s", err)
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

func DeleteStoryRequest(c *fiber.Ctx, db *sql.DB) error {
	context := c.Context()
	headers := c.GetReqHeaders()
	storyId := headers["Id"]
	err := queries.DeleteSingleStory(storyId, context, db)
	if err != nil {
		log.Fatalf("Delete Story Fatal Error - %s", err)
	}

	message := fmt.Sprintf("Deleted story with id: %s", storyId)
	response := models.JSONResponse{
		Type:    "success",
		Data:    nil,
		Message: message,
	}

	c.Response().StatusCode()
	c.Response().Header.Add("Content-Type", "application/json")
	return c.JSON(response)
}

func UpdateStoryRequest(c *fiber.Ctx, db *sql.DB) error {
	queryError := queries.UpdateStory(c, db)
	if queryError != nil {
		return response.BasicErrorHandling(c, queryError)
	}

	headers := c.GetReqHeaders()
	id := headers["Id"]
	message := fmt.Sprintf("Updated story with id: %s", id)
	response := models.JSONResponse{
		Type:    "success",
		Message: message,
		Data:    nil,
	}

	c.Response().StatusCode()
	c.Response().Header.Add("Content-Type", "application/json")
	return c.JSON(response)
}
