package dynamodbrequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
	configuration "github.com/tolubydesign/angular-story-backend/app/config"
	database "github.com/tolubydesign/angular-story-backend/app/database"
	"github.com/tolubydesign/angular-story-backend/app/handler"
	helpers "github.com/tolubydesign/angular-story-backend/app/helpers"
	"github.com/tolubydesign/angular-story-backend/app/logging"
	models "github.com/tolubydesign/angular-story-backend/app/models"
	"github.com/tolubydesign/angular-story-backend/app/mutation"
	query "github.com/tolubydesign/angular-story-backend/app/query"
	"github.com/tolubydesign/angular-story-backend/app/utils"
)

// Get client connection to dynamodb.
func ConnectDynamoDB() (*dynamodb.Client, error) {
	// Connect to PostgreSQL database
	client, err := database.CreateDynamoClient()
	if err != nil {
		return nil, err
	}

	return client, nil
}

/*
Function incomplete.
Creating a new table.
Database name and structure are provided by user request and request body.
*/
func CreateNewTable(ctx *fiber.Ctx, client *dynamodb.Client) error {
	if client == nil {
		return errors.New("Dynamo Database inaccessible.")
	}

	// TODO: get table name from request context. For now use hardcoded value
	// TODO: get function parameters from request context
	// TODO: use table.CreateDynamoDBTable to dynamically create new tables, when needed

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(models.HTTPResponse{
		Code:    fiber.StatusForbidden,
		Message: "Function has not been created",
		Error:   true,
	})
}

func GetAllDynamoDBTables(ctx *fiber.Ctx, client *dynamodb.Client, c *configuration.Config) error {
	log.Println("Getting all tables within the dynamo database")
	var err error
	tableName := c.Configuration.Dynamodb.StoryTableName
	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}

	tables, err := table.ListDynamodbTables()
	if err != nil {
		// Return error to user
		return fiber.NewError(fiber.StatusInternalServerError, error.Error(err))
	}

	response := models.HTTPResponse{
		Code:    fiber.StatusOK,
		Data:    tables,
		Message: "Tables found in DynamoDB.",
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

/*
Generate false data to populate dynamodb database
*/
func PopulateDynamoDatabase(ctx *fiber.Ctx, client *dynamodb.Client, c *configuration.Config) error {
	logging.Event("Adding default data to the dynamo database")
	var err error

	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      c.Configuration.Dynamodb.StoryTableName,
	}

	err = helpers.PopulateStoryDatabase(table)
	if err != nil {
		// Return error to user
		return fiber.NewError(fiber.StatusInternalServerError, error.Error(err))
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(models.HTTPResponse{
		Code:    fiber.StatusOK,
		Message: "Request successful",
	})
}

/*
List all stories within the database.

- https://towardsdatascience.com/dynamodb-go-sdk-how-to-use-the-scan-and-batch-operations-efficiently-5b41988b4988
*/
func ListAllStoriesRequest(ctx *fiber.Ctx, client *dynamodb.Client, c *configuration.Config) error {
	logging.Event("Getting all stories request.")
	if client == nil {
		return errors.New(helpers.ResponseMessages.NilClient)
	}

	tableName := c.Configuration.Dynamodb.StoryTableName
	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}

	items, err := table.FullTableScan()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := models.HTTPResponse{
		Code:    fiber.StatusOK,
		Data:    items,
		Message: "Request successful",
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

/*
REQUEST. Add new story to dynamodb. Content for story is required.
*/
func AddStoryRequest(ctx *fiber.Ctx, client *dynamodb.Client, c *configuration.Config) error {
	logging.Event("[Enter]: Adding Story request")
	var err error
	if client == nil {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusInternalServerError,
			Message: helpers.ResponseMessages.NilClient,
		})
	}

	id, err := AddStory(ctx, client, c)
	if err != nil {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	logging.Event("[Exit]: Adding Story request")
	ctx.Response().StatusCode()
	return handler.HandleResponse(handler.ResponseHandlerParameters{
		Context: ctx,
		Code:    fiber.StatusOK,
		Message: fmt.Sprint("Content added id: ", id),
	})
}

/*
REQUEST Get Story bases on id provided in request.

Will error if no id header is found
*/
func GetStoryByIdRequest(ctx *fiber.Ctx, client *dynamodb.Client, c *configuration.Config) error {
	var err error
	var story *models.DynamoStoryResponseStruct

	logging.Event("[Enter] Getting Single Story")
	if client == nil {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusInternalServerError,
			Message: helpers.ResponseMessages.NilClient,
		})
	}

	storyId, err := utils.GetRequestHeaderID(ctx)
	if (len(storyId) < 6) || (storyId == "") {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Make sure that the provided id is valid
	if (len(storyId) < 6) || (storyId == "") {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ID provided",
		})
	}

	logging.Event("Found story: ", storyId)

	// Setup table
	tableName := c.Configuration.Dynamodb.StoryTableName
	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}

	// Get story based on id provided in request
	story, err = table.GetStoryById(storyId)
	if err != nil {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// For when the database was scanned but nothing was found
	if story == nil {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusNotFound,
			Message: "Not found",
		})
	}

	logging.Event("[Exit] Getting Single Story")
	return handler.HandleResponse(handler.ResponseHandlerParameters{
		Context: ctx,
		Code:    fiber.StatusOK,
		Data:    story,
		Message: "Request successfully",
	})
}

/*
Update story based on the story id and body context provided.

If the user does not provide both id and context, with "title", "description" and "content" (in the request body),
the request will return an error.
*/
func UpdateDynamodbStoryRequest(ctx *fiber.Ctx, client *dynamodb.Client, c *configuration.Config) error {
	logging.Event("Update dynamodb story request.")
	if client == nil {
		logging.Error("Client is null.")
		return fiber.NewError(fiber.StatusInternalServerError, helpers.ResponseMessages.NilClient)
	}

	content, fiberError := DynamoUpdateStory(ctx, client, c)
	if fiberError != nil {
		return fiberError
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(models.HTTPResponse{
		Code:    fiber.StatusOK,
		Data:    models.DynamoStoryResponseStruct(content),
		Message: "Request successful",
	})
}

/*
Delete single story item in dynamo database.

Returning fiber error if issues occur
*/
func DeleteDynamodbStoryRequest(ctx *fiber.Ctx, client *dynamodb.Client, c *configuration.Config) error {
	logging.Event("[Enter] Delete dynamodb story request")
	if client == nil {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusInternalServerError,
			Message: helpers.ResponseMessages.NilClient,
		})
		// return fiber.NewError(fiber.StatusInternalServerError, helpers.ResponseMessages.NilClient)
	}

	id, err := utils.GetRequestHeaderID(ctx)
	if err != nil {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
		// return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	creator := helpers.GetRequestHeader(ctx, "Creator")
	log.Println("Delete dynamodb story request. creator ", creator)
	if creator == "" {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request. Creator ID required",
		})
	}

	// Verify that id is a valid uuid
	v := utils.ValidUuid(id)
	if v != true {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusBadRequest,
			Message: "Invalid id provided",
		})
	}

	id, err = DynamoDeleteStory(ctx, client, id, creator, c)
	if err != nil {
		return handler.HandleResponse(handler.ResponseHandlerParameters{
			Context: ctx,
			Error:   true,
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return handler.HandleResponse(handler.ResponseHandlerParameters{
		Context: ctx,
		Code:    fiber.StatusOK,
		Message: fmt.Sprintf("Delete story id: %s", id),
	})
}

/*
List all users in dynamodb database
*/
func ListAllUsersRequest(ctx *fiber.Ctx, client *dynamodb.Client, c *configuration.Config) error {
	var userJson []models.User
	if client == nil {
		return errors.New(helpers.ResponseMessages.NilClient)
	}

	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      c.Configuration.Dynamodb.UserTableName,
	}

	items, err := table.FullUserTableScan()
	if err != nil {
		// Return error to user
		return fiber.NewError(fiber.StatusInternalServerError, error.Error(err))
	}

	// Convert items ([]models.DatabaseUserStruct) into json
	temp, _ := json.Marshal(items)
	err = json.Unmarshal(temp, &userJson)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(models.HTTPResponse{
		Code:    fiber.StatusOK,
		Data:    userJson,
		Message: "Request successful",
	})
}

// Login (Sign In) to user account
func UserLoginRequest(c *fiber.Ctx, client *dynamodb.Client) error {
	var user models.DatabaseUserStruct
	if client == nil {
		return errors.New(helpers.ResponseMessages.NilClient)
	}

	// Get user login information. Email and password
	user, err := helpers.GetLoginInfoContext(c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	verified, err := UserVerification(client, user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Response().StatusCode()
	c.Response().Header.Add("Content-Type", "application/json")
	return c.JSON(models.HTTPResponse{
		Code:    fiber.StatusOK,
		Data:    verified,
		Message: "User login successful",
	})
}

// Sign Up user to database
func UserSignUpRequest(c *fiber.Ctx, client *dynamodb.Client, config *configuration.Config) error {
	if client == nil {
		return errors.New(helpers.ResponseMessages.NilClient)
	}

	err := SignUpUser(c, client, config)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Response().StatusCode()
	c.Response().Header.Add("Content-Type", "application/json")
	return c.JSON(models.HTTPResponse{
		Code:    fiber.StatusOK,
		Message: "User Signed Up successful",
	})
}

func HealthCheck(ctx *fiber.Ctx, client *dynamodb.Client, c *configuration.Config) error {
	var response models.HTTPResponse
	if client == nil {
		return errors.New(helpers.ResponseMessages.NilClient)
	}

	logging.Event("Performing health check on dynamodb.")
	tableName := c.Configuration.Dynamodb.StoryTableName
	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}

	tables, err := table.ListDynamodbTables()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response = models.HTTPResponse{
		Code:    fiber.StatusOK,
		Message: "Dynamo Database table unreachable.",
	}

	// Check for story table in array of table names
	for _, a := range tables {
		if a == tableName {
			response = models.HTTPResponse{
				Code:    fiber.StatusOK,
				Message: "DynamoDB is active.",
			}
		}
	}

	logging.Event("Performing health check on dynamodb completed.")
	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}
