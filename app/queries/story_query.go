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
	"github.com/valyala/fasthttp"
)

type DatabaseService struct {
	db *sql.DB
}

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

func GetSingleStory(id string, con *fasthttp.RequestCtx, db *sql.DB) (models.Story, error) {
	ctx, cancel := context.WithTimeout(con, 6*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT * FROM story
		WHERE story_id = '%s';
	`, id)
	var story models.Story

	err := db.QueryRowContext(ctx, query).Scan(&story.StoryId, &story.Title, &story.Description, &story.Content)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("no story with id of %s\n", id)
	case err != nil:
		log.Fatalf("query error: %v\n", err)
	default:
		log.Printf("story with id of %s found.\n", id)
	}

	return story, nil
}

func AddStory(c *fiber.Ctx, db *sql.DB) error {
	fiberContext := c.Context()
	ctx, cancel := context.WithTimeout(fiberContext, 3*time.Second)
	defer cancel()

	var body models.Story
	byteBody := c.Body()
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
		// log.Fatalf("Fatal Results Error - %s", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
		// log.Fatalf("Fatal Rows Error: %d", err)
	}
	if rows != 1 {
		return err
		// log.Fatalf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}

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
