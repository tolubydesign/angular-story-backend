package queries

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

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
