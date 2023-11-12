package models

import (
	"github.com/google/uuid"
)

type Story struct {
	Id          uuid.UUID   `json:"id" validate:"uuid"`
	Title       string      `json:"title" validate:"required"`
	Description string      `json:"description" validate:"required"`
	Content     interface{} `json:"content"` // type *StoryContent
}

type StoryContent struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Children    *[]StoryContent `json:"children"`
}

type DynamoStoryResponseStruct struct {
	Id          string        `json:"id" validate:"uuid"`
	Creator     string        `json:"creator" validate:"required"`
	Title       string        `json:"title" validate:"required"`
	Description string        `json:"description" validate:"required"`
	Content     *StoryContent `json:"content"`
}

type DynamoStoryDatabaseStruct struct {
	// Dynamodb key (id & creator)
	// creator value must connect to user.id or "default"
	Id          string        `dynamodbav:"id"`
	Creator     string        `dynamodbav:"creator"`
	Title       string        `dynamodbav:"title"`
	Description string        `dynamodbav:"description"`
	Content     *StoryContent `dynamodbav:"content"`
}
