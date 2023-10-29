package controller

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
	database "github.com/tolubydesign/angular-story-backend/app/database"
	helpers "github.com/tolubydesign/angular-story-backend/app/helpers"
	models "github.com/tolubydesign/angular-story-backend/app/models"
	mutation "github.com/tolubydesign/angular-story-backend/app/mutation"
	query "github.com/tolubydesign/angular-story-backend/app/query"
)

/*
Get client connection to dynamodb.
*/
func ConnectDynamoDB() (*dynamodb.Client, error) {
	// Connect to PostgreSQL database
	client, err := database.CreateDynamoClient()
	if err != nil {
		return nil, err
	}

	return client, nil
}

/*
Creating a new table
TODO: write description for function
{...}
*/
func CreateNewTable(ctx *fiber.Ctx, client *dynamodb.Client) error {
	if client == nil {
		return errors.New("Dynamo Database inaccessible.")
	}

	// TODO: get table name from request context. For now use hardcoded value
	// TODO: get function parameters from request context
	// TODO: use table.CreateDynamoDBTable to dynamically create new tables, when needed
	response := models.JSONResponse{
		Type:    "mid",
		Data:    nil,
		Message: "Function has not been created",
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

func GetAllDynamoDBTables(ctx *fiber.Ctx, client *dynamodb.Client) error {
	fmt.Printf("Getting all tables within the dynamo database")

	tableName := "Story"
	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}

	tables, err := table.ListDynamodbTables()

	if err != nil {
		// Return error to user
		return fiber.NewError(fiber.StatusInternalServerError, error.Error(err))
	}

	response := models.JSONResponse{
		Type:    "success",
		Data:    tables,
		Message: "Tables found in DynamoDB.",
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

func PopulateDynamoDatabase(ctx *fiber.Ctx, client *dynamodb.Client) error {
	fmt.Println("Adding default data to the dynamo database")
	tableName := "Story"
	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}

	err := helpers.PopulateStoryDatabase(table)
	if err != nil {
		// Return error to user
		return fiber.NewError(fiber.StatusInternalServerError, error.Error(err))
	}

	response := models.JSONResponse{
		Type:    "success",
		Data:    nil,
		Message: "Request successful",
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

/*
TODO: add description

param - ctx Fiber Context

param - client Dynamo DB client

@see - https://towardsdatascience.com/dynamodb-go-sdk-how-to-use-the-scan-and-batch-operations-efficiently-5b41988b4988
{...}
*/
func ListAllStories(ctx *fiber.Ctx, client *dynamodb.Client) error {
	if client == nil {
		return errors.New(helpers.DynamodbResponseMessages["nilClient"])
	}

	tableName := "Story"
	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}

	items, err := table.FullTableScan()

	if err != nil {
		// Return error to user
		return fiber.NewError(fiber.StatusInternalServerError, error.Error(err))
	}

	response := models.JSONResponse{
		Type:    "success",
		Data:    items,
		Message: "Request successful",
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

/*
Add new story to dynamodb

ctx - Fiber Context

client - Dynamo DB client
*/
func AddStoryRequest(ctx *fiber.Ctx, client *dynamodb.Client) error {
	fmt.Println("Adding Story request.")

	var err error
	if client == nil {
		return fiber.NewError(fiber.StatusInternalServerError, helpers.DynamodbResponseMessages["nilClient"])
	}

	// Setup table
	tableName := "Story"
	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}

	// Get data from fiber context
	var body models.DynamoStoryDatabaseStruct
	byteBody := ctx.Body()

	// Convert Struct to JSON
	json.Unmarshal(byteBody, &body)
	// json, err := json.Marshal(body.Content)
	if err != nil {
		return err
	}

	// Convert Struct to JSON
	json.Unmarshal(byteBody, &body)
	if err != nil {
		return err
	}

	story := models.DynamoStoryDatabaseStruct{
		Id:          helpers.GenerateStringUUID(),
		Title:       body.Title,
		Description: body.Description,
		Content:     body.Content,
	}

	err = table.AddStory(story)

	if err != nil {
		// TODO: create better http response
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := models.JSONResponse{
		Type:    "success",
		Data:    nil,
		Message: "Story added successfully",
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

/*
{...}
*/
func GetStoryByIdRequest(ctx *fiber.Ctx, client *dynamodb.Client) error {
	// TODO: log event
	fmt.Println("Get story by id request.")

	var err error
	var story *models.DynamoStoryResponseStruct
	if client == nil {
		return fiber.NewError(fiber.StatusInternalServerError, helpers.DynamodbResponseMessages["nilClient"])
	}

	headers := ctx.GetReqHeaders()
	storyId := headers["Id"]

	fmt.Println("Story id", storyId)
	if (len(storyId) < 6) || (storyId == "") {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid ID provided")
	}

	// Setup table
	tableName := "Story"
	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}

	// Get story id, provided in request
	story, err = table.GetStoryById(storyId)
	if err != nil {
		// TODO: create better http response
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(models.JSONResponse{
		Type:    "success",
		Data:    story,
		Message: "Request successfully",
	})
}
