package models

import (
	"github.com/google/uuid"
)

type StoryContent struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Children    interface{} `json:"children"`
}
type Story struct {
	StoryId     uuid.UUID   `json:"story_id" validate:"uuid"`
	Title       string      `json:"title" validate:"required"`
	Description string      `json:"description" validate:"required"`
	Content     interface{} `json:"content"`
}

type AllStories struct {
	Story []Story `json:"story"`
}
