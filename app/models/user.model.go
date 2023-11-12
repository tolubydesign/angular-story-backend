package models

// User JSON structure
type User struct {
	Id           string `json:"id" validate:"uuid"`
	Email        string `json:"email" validate:"required"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	AccountLevel string `json:"accountLevel" validate:"required"`
}

// User database structure
type DatabaseUserStruct struct {
	// Dynamodb key (id | email)
	// id connected to story.creator
	Id           string `dynamodbav:"id"`
	Email        string `dynamodbav:"email"`
	Name         string `dynamodbav:"name"`
	Surname      string `dynamodbav:"surname"`
	Username     string `dynamodbav:"username"`
	Password     string `dynamodbav:"password"`
	AccountLevel string `dynamodbav:"accountLevel"`
}
