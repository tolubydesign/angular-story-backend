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
	Creator     string        `json:"creator" validate:"required"`
	Title       string        `json:"title" validate:"required"`
	Description string        `json:"description" validate:"required"`
	Content     *StoryContent `json:"content"`
}

type DynamoStoryDatabaseStruct struct {
	// Dynamodb key
	Id string `dynamodbav:"id"`
	// Dynamodb key
	// creator value must connect to user.id or "default"
	Creator     string        `dynamodbav:"creator"`
	Title       string        `dynamodbav:"title"`
	Description string        `dynamodbav:"description"`
	Content     *StoryContent `dynamodbav:"content"`
}

// User structure

// User JSON structure
type User struct {
	Id           string  `json:"id" validate:"uuid"`
	Email        string  `json:"email" validate:"required"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Username     string  `json:"username" validate:"required"`
	Password     string  `json:"password" validate:"required"`
	AccountLevel string  `json:"accountLevel" validate:"required"`
}

// User database structure
type DatabaseUserStruct struct {
	// Dynamodb key
	Id string `dynamodbav:"id" validate:"uuid"`
	// Dynamodb key
	Email        string  `dynamodbav:"email" validate:"required"`
	Username     string  `dynamodbav:"username" validate:"required"`
	Name         string `dynamodbav:"name"`
	Surname      string `dynamodbav:"surname"`
	Password     string  `dynamodbav:"password" validate:"required"`
	AccountLevel string  `dynamodbav:"accountLevel"`
}
