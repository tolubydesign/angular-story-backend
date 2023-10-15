package models

import (
	"github.com/google/uuid"
)

type StoryContent struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Children    *[]StoryContent `json:"children"`
}
type Story struct {
	Id          uuid.UUID   `json:"id" validate:"uuid"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	Content     interface{} `json:"content"` // type *StoryContent
}

type AllStories struct {
	Story []Story `json:"story"`
}
