package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"

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
