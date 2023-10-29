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

type AllStories struct {
	Story []Story `json:"story"`
}

type DynamoStoryResponseStruct struct {
	Id          string        `json:"id" validate:"uuid"`
	Title       string        `json:"title" validate:"required"`
	Description string        `json:"description" validate:"required"`
	Content     *StoryContent `json:"content"`
}

type DynamoStoryDatabaseStruct struct {
	Id          string        `dynamodbav:"id"`
	Title       string        `dynamodbav:"title"`
	Description string        `dynamodbav:"description"`
	Content     *StoryContent `dynamodbav:"content"`
}

type User struct {
	Id       string  `json:"id" validate:"uuid"`
	Username string  `json:"username" validate:"required"`
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Email    string  `json:"email" validate:"required"`
	Password string  `json:"password" validate:"required"`
}

type DatabaseUserStruct struct {
	Id       string  `dynamodbav:"id" validate:"uuid"`
	Username string  `dynamodbav:"username" validate:"required"`
	Name     *string `dynamodbav:"name"`
	Surname  *string `dynamodbav:"surname"`
	Email    string  `dynamodbav:"email" validate:"required"`
	Password string  `dynamodbav:"password" validate:"required"`
}
