package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	database "github.com/tolubydesign/angular-story-backend/app/database"
	"github.com/tolubydesign/angular-story-backend/app/models"
	query "github.com/tolubydesign/angular-story-backend/app/query"
	"github.com/tolubydesign/angular-story-backend/app/utils"
	"github.com/tolubydesign/angular-story-backend/pkg/response"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

/*
GET request. Get all stories in database.

Return database response or possible Error.
*/
func GetAllStoriesRequest(ctx *fiber.Ctx, db *sql.DB) error {
	var error error
	error = database.LogEvent("REQUEST START: all stories")
	if error != nil {
		// TODO: create better http response
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	stories, err := query.GetAllStories(db)
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
		// TODO Handle error response
		panic(err)
	}

	error = database.LogEvent("REQUEST SUCCESSFUL: all stories")
	if error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

/*
GET request. Get a single story based on the id provided in the request header.

Returning fiber context response or possible Error.
*/
func GetSingleStoryRequest(ctx *fiber.Ctx, db *sql.DB) error {
	context := ctx.Context()
	headers := ctx.GetReqHeaders()
	headerId := headers["Id"]

	var error error

	// logging
	error = database.LogEvent(fmt.Sprintf("REQUEST START: single story, id:%s", headerId))
	if error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	story, error := query.GetSingleStory(headerId, context, db)
	if error != nil {
		panic(error)
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

	error = database.LogEvent(fmt.Sprintf("REQUEST SUCCESSFUL: single story, id:%s", headerId))
	if error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

/*
Confirm the status of the database.

Returning fiber context response or possible Error.
*/
func CheckHealthRequest(c *fiber.Ctx, db *sql.DB) error {
	var error error
	error = database.LogEvent("REQUEST START: health check")
	if error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	ctxt := c.Context()
	ctx, cancel := context.WithTimeout(ctxt, 2*time.Second)

	defer cancel()
	err := db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	error = database.LogEvent("REQUEST SUCCESSFUL: health check")
	if error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	c.Response().StatusCode()
	return c.SendString("Database is working.")
}

/*
POST request. Creates new story.

Returning fiber context response or possible Error.
*/
func InsertStoryRequest(c *fiber.Ctx, db *sql.DB) error {
	var error error
	error = database.LogEvent("REQUEST START: insert story")
	if error != nil {
		// TODO: create better http response
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	// TODO: verify that information being sent in is valid json
	err := query.AddStory(c, db)
	if err != nil {
		log.Fatalf("Insert Story Fatal Error - %s", err)
	}

	response := models.JSONResponse{
		Type:    "success",
		Data:    nil,
		Message: "Database has been updated.",
	}

	error = database.LogEvent("REQUEST SUCCESSFUL: insert story")
	if error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	c.Response().StatusCode()
	c.Response().Header.Add("Content-Type", "application/json")
	return c.JSON(response)
}

/*
DELETE request. Delete story based on id provided.

Returning fiber context response or possible Error.
*/
func DeleteStoryRequest(c *fiber.Ctx, db *sql.DB) error {
	context := c.Context()
	id := utils.GetRequestHeaderId(c)

	var error error
	error = database.LogEvent(fmt.Sprintf("REQUEST START: delete story, id:%s", id))
	if error != nil {
		// TODO: create better http response
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	err := query.DeleteSingleStory(id, context, db)
	if err != nil {
		log.Fatalf("Delete Story Fatal Error - %s", err)
	}

	message := fmt.Sprintf("Deleted story with id: %s", id)
	response := models.JSONResponse{
		Type:    "success",
		Data:    nil,
		Message: message,
	}

	error = database.LogEvent(fmt.Sprintf("REQUEST SUCCESSFUL: delete story, id:%s", id))
	if error != nil {
		// TODO: create better http response
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	c.Response().StatusCode()
	c.Response().Header.Add("Content-Type", "application/json")
	return c.JSON(response)
}

/*
PUT request. Update existing story information, based on header id provided.

Returning fiber context response or possible Error.
*/
func UpdateStoryRequest(c *fiber.Ctx, db *sql.DB) error {
	id := utils.GetRequestHeaderId(c)

	var error error
	error = database.LogEvent(fmt.Sprintf("REQUEST START: update story, id:%s", id))
	if error != nil {
		// TODO: create better http response
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	queryError := query.UpdateStory(c, db)
	if queryError != nil {
		return response.BasicErrorHandling(c, queryError)
	}

	message := fmt.Sprintf("Updated story with id: %s", id)
	response := models.JSONResponse{
		Type:    "success",
		Message: message,
		Data:    nil,
	}

	error = database.LogEvent(fmt.Sprintf("REQUEST SUCCESSFUL: update story, id:%s", id))
	if error != nil {
		// TODO: create better http response
		return fiber.NewError(fiber.StatusInternalServerError, error.Error())
	}

	c.Response().StatusCode()
	c.Response().Header.Add("Content-Type", "application/json")
	return c.JSON(response)
}
