package queries

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/tolubydesign/angular-story-backend/app/models"
	"github.com/tolubydesign/angular-story-backend/app/utils"
	"github.com/valyala/fasthttp"
)

func GetAllStories(db *sql.DB) ([]models.Story, error) {
	request := `select * from Story`
	rows, err := db.Query(request)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stories []models.Story

	for rows.Next() {
		var story models.Story
		err := rows.Scan(&story.Id, &story.Title, &story.Description, &story.Content)
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

func GetSingleStory(id string, con *fasthttp.RequestCtx, db *sql.DB) (models.Story, error) {
	ctx, cancel := context.WithTimeout(con, 6*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT * FROM story
		WHERE story_id = '%s';
	`, id)
	var story models.Story
	err := db.QueryRowContext(ctx, query).Scan(&story.Id, &story.Title, &story.Description, &story.Content)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("Get Single Story. Error no story with id of %s\n", id)
	case err != nil:
		log.Fatalf("Get Single Story. query error: %v\n", err)
	default:
		log.Printf("Get Single Story. story with id of %s found.\n", id)
	}

	if err != nil {
		return story, err
	}

	return story, nil
}

// POST Request.
// Add a new story to the database. Includes both content and
func AddStory(c *fiber.Ctx, db *sql.DB) error {
	fiberContext := c.Context()
	ctx, cancel := context.WithTimeout(fiberContext, 3*time.Second)
	defer cancel()

	var body models.Story
	byteBody := c.Body()

	// Convert Struct to JSON
	json.Unmarshal(byteBody, &body)
	bodyContentJson, err := json.Marshal(body.Content)
	if err != nil {
		return err
	}

	model := models.Story{
		Title:       body.Title,
		Description: body.Description,
		Content:     bodyContentJson,
	}

	execution := "INSERT INTO story(title, description, content) VALUES($1, $2, $3);"
	result, err := db.ExecContext(ctx, execution, model.Title, model.Description, model.Content)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return err
	}

	return nil
}

// DELETE Request.
// Remove a single story.
func DeleteSingleStory(id string, ctx *fasthttp.RequestCtx, db *sql.DB) error {
	basicContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	deleteStmt := `DELETE FROM story WHERE story_id=$1`
	_, err := db.ExecContext(basicContext, deleteStmt, id)
	if err != nil {
		return err
	}

	return nil
}

// PUT Request.
// Update the contents and information of a single story row.
func UpdateStory(ctx *fiber.Ctx, db *sql.DB) error {
	fiberContext := ctx.Context()
	headers := ctx.GetReqHeaders()
	headerId := headers["Id"]
	headerDescription := headers["Description"]
	headerTitle := headers["Title"]

	basicContext, cancel := context.WithTimeout(fiberContext, 2*time.Second)
	defer cancel()

	var body models.Story
	byteBody := ctx.Body()
	json.Unmarshal(byteBody, &body)
	content, err := json.Marshal(body.Content)
	if err != nil {
		return err
	}

	headerIdError := utils.ValidateLimitedStringVariable(headerId)
	if headerIdError != nil {
		return headerIdError
	}

	headerTitleError := utils.ValidateLimitedStringVariable(headerTitle)
	if headerTitleError != nil {
		return headerTitleError
	}

	headerDescriptionError := utils.ValidateLimitedStringVariable(headerDescription)
	if headerDescriptionError != nil {
		return headerDescriptionError
	}

	var updateStmt string
	updateStmt = fmt.Sprintf(`
	UPDATE story 
	SET	title = $1,
			description = $2,
			content = $3
	WHERE story_id = $4;
	`)

	result, err := db.ExecContext(basicContext, updateStmt, headerTitle, headerDescription, content, headerId)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return err
	}

	return nil
}
