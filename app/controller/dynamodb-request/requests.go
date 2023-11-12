package dynamodbrequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
	"github.com/tolubydesign/angular-story-backend/app/config"
	database "github.com/tolubydesign/angular-story-backend/app/database"
	helpers "github.com/tolubydesign/angular-story-backend/app/helpers"
	models "github.com/tolubydesign/angular-story-backend/app/models"
	mutation "github.com/tolubydesign/angular-story-backend/app/mutation"
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
func PopulateDynamoDatabase(ctx *fiber.Ctx, client *dynamodb.Client) error {
	fmt.Println("Adding default data to the dynamo database")
	configuration, err := config.GetConfiguration()
	if err != nil {
		// Return error to user
		return fiber.NewError(fiber.StatusInternalServerError, error.Error(err))
	}

	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      configuration.Configuration.Dynamodb.StoryTableName,
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
func ListAllStoriesRequest(ctx *fiber.Ctx, client *dynamodb.Client) error {
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
Add new story to dynamodb. Content for story is required
*/
func AddStoryRequest(ctx *fiber.Ctx, client *dynamodb.Client) error {
	log.Println("Adding Story request.")
	var err error
	configuration, err := config.GetConfiguration()

	if client == nil {
		return fiber.NewError(fiber.StatusInternalServerError, helpers.DynamodbResponseMessages["nilClient"])
	}

	// Setup table
	// TODO: get table name from env
	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      configuration.Configuration.Dynamodb.StoryTableName,
	}

	story, err := helpers.GenerateStoryFromRequestContext(ctx)
	if err != nil {
		return err
	}

	// Generate new id for story
	story.Id = helpers.GenerateStringUUID()
	err = table.AddStory(story)
	if err != nil {
		// TODO: create better http response
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := models.HTTPResponse{
		Code:    fiber.StatusOK,
		Message: "Story added successfully",
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(response)
}

/*
Get Story bases on id provided in request.

Will error if no id header is found
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
		// TODO: more descriptive response
		ctx.Response().StatusCode()
		ctx.Response().Header.Add("Content-Type", "application/json")
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
	return ctx.JSON(models.HTTPResponse{
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
func UpdateDynamodbStoryRequest(ctx *fiber.Ctx, client *dynamodb.Client) error {
	log.Println("Update dynamodb story request.")
	var err error
	configuration, err := config.GetConfiguration()
	if client == nil {
		return errors.New(helpers.DynamodbResponseMessages["nilClient"])
	}

	id, err := helpers.GetRequestHeaderID(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	creator := helpers.GetRequestHeader(ctx, "Creator")
	validCreator := utils.ValidateString(creator)
	if validCreator != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid creator value provided.")
	}

	// Get body context
	// TODO: verify structure of body json provided
	story, err := helpers.GenerateStoryFromRequestContext(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      configuration.Configuration.Dynamodb.StoryTableName,
	}

	story = models.DynamoStoryDatabaseStruct{
		Id:          id,
		Creator:     creator,
		Title:       story.Title,
		Description: story.Description,
		Content:     story.Content,
	}

	// Update story, in database, from content provided.
	content, err := table.UpdateDynamoDBTable(story)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// TODO: add undefined parameters, content.creator & content.id, in response.

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
*/
func DeleteDynamodbStoryRequest(ctx *fiber.Ctx, client *dynamodb.Client) error {
	log.Println("Delete dynamodb story request.")
	var err error
	configuration, err := config.GetConfiguration()
	if client == nil {
		return errors.New(helpers.DynamodbResponseMessages["nilClient"])
	}

	id, err := helpers.GetRequestHeaderID(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	creator := helpers.GetRequestHeader(ctx, "Creator")
	log.Println("Delete dynamodb story request. creator ", creator)
	if creator == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Story Title not provided.")
	}

	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      configuration.Configuration.Dynamodb.StoryTableName,
	}

	story := models.DynamoStoryDatabaseStruct{
		Id:      id,
		Creator: creator,
	}

	// Update story, in database, from content provided.
	err = table.DeleteStory(story)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	ctx.Response().StatusCode()
	ctx.Response().Header.Add("Content-Type", "application/json")
	return ctx.JSON(models.HTTPResponse{
		Code:    fiber.StatusOK,
		Message: fmt.Sprintf("Request to delete story successful id:%s", id),
	})
}

/*
List all users in dynamodb database
*/
func ListAllUsersRequest(ctx *fiber.Ctx, client *dynamodb.Client) error {
	var userJson []models.User
	configuration, err := config.GetConfiguration()
	if client == nil {
		return errors.New(helpers.DynamodbResponseMessages["nilClient"])
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      configuration.Configuration.Dynamodb.UserTableName,
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
		return errors.New(helpers.DynamodbResponseMessages["nilClient"])
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
func UserSignUpRequest(c *fiber.Ctx, client *dynamodb.Client) error {
	if client == nil {
		return errors.New(helpers.DynamodbResponseMessages["nilClient"])
	}

	err := SignUpUser(c, client)
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
