package dynamodbrequest

import (
	"errors"

	// "fmt"
	"log"

	// "log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
	"github.com/tolubydesign/angular-story-backend/app/config"
	configuration "github.com/tolubydesign/angular-story-backend/app/config"
	"github.com/tolubydesign/angular-story-backend/app/helpers"
	"github.com/tolubydesign/angular-story-backend/app/models"
	"github.com/tolubydesign/angular-story-backend/app/mutation"
	"github.com/tolubydesign/angular-story-backend/app/query"
	"github.com/tolubydesign/angular-story-backend/app/utils"
)

/*
Check if user credentials is a valid user. Find user via email and name
*/
func UserVerification(client *dynamodb.Client, user models.DatabaseUserStruct) (string, error) {
	var token string
	var err error

	configuration, err := config.GetConfiguration()
	if err != nil {
		return token, err
	}

	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      configuration.Configuration.Dynamodb.UserTableName,
	}

	u, err := table.GetUserByEmail(user.Email)
	if err != nil {
		return token, err
	}

	// Validate that password provided matches known (hashed) password, in database
	// If user is found. Check if password matches database password
	if u.Email != user.Email || u.Password != user.Password {
		return token, errors.New("User information doesn't match known user information.")
	}

	// Set Token
	token, err = utils.BuildUserJWTToken(u)
	if err != nil {
		return token, err
	}

	return token, err
}

func SignUpUser(c *fiber.Ctx, client *dynamodb.Client, config *configuration.Config) error {
	u, err := helpers.GetSignInInfoContext(c)
	if err != nil {
		return err
	}

	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      config.Configuration.Dynamodb.UserTableName,
	}

	// Hash password
	hash := u.Password

	err = table.AddUser(models.DatabaseUserStruct{
		Id:           helpers.GenerateStringUUID(),
		Email:        u.Email,
		Username:     u.Username,
		Name:         u.Name,
		Surname:      u.Surname,
		Password:     hash,
		AccountLevel: u.AccountLevel,
	})

	if err != nil {
		return err
	}

	return nil
}

// FUNCTION RELATED REQUEST. Make a DynamoDB request to
func DynamoDeleteStory(c *fiber.Ctx, client *dynamodb.Client, id string, creator string, config *configuration.Config) (string, error) {
	var err error
	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      config.Configuration.Dynamodb.StoryTableName,
	}

	story := models.DynamoStoryDatabaseStruct{
		Id:      id,
		Creator: creator,
	}

	// Update story, in database, from content provided.
	err = table.DeleteStory(story)
	return id, err
}

// FUNCTION RELATED REQUEST. Make request to Dynamodb. Update a single story. Including title, description and content.
func DynamoUpdateStory(c *fiber.Ctx, client *dynamodb.Client, config *config.Config) (models.DynamoStoryDatabaseStruct, *fiber.Error) {
	var content models.DynamoStoryDatabaseStruct
	var err error

	// Get story id
	// TODO: move this value to request body
	id, err := utils.GetRequestHeaderID(c)
	if err != nil {
		return content, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Verify that id provided is valid
	v := utils.ValidUuid(id)
	if v != true {
		return content, fiber.NewError(fiber.StatusBadRequest, helpers.ResponseMessages.InvalidUUID)
	}

	// Get id of owner of story
	// TODO: move this value to request body
	creator := helpers.GetRequestHeader(c, "Creator")
	validCreator := utils.ValidateString(creator)
	if validCreator != nil {
		return content, fiber.NewError(fiber.StatusBadRequest, helpers.ResponseMessages.InvalidCreatorID)
	}

	// Get body context
	// TODO: verify structure of body json provided
	story, err := helpers.GetStoryFromRequestContext(c)
	if err != nil {
		return content, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      config.Configuration.Dynamodb.StoryTableName,
	}

	story = models.DynamoStoryDatabaseStruct{
		Id:          id,
		Creator:     creator,
		Title:       story.Title,
		Description: story.Description,
		Content:     story.Content,
	}

	// Update story, in database, from content provided.
	content, err = table.UpdateDynamoDBTable(story)
	if err != nil {
		return content, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// TODO: add undefined parameters, content.creator & content.id, in response.
	return content, nil
}

// FUNCTION RELATED REQUEST. Add a single story to database.
// Return basic error is issues occur
func AddStory(c *fiber.Ctx, client *dynamodb.Client, config *configuration.Config) (string, error) {
	log.Println("Executing Adding Story function")
	var id string

	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      config.Configuration.Dynamodb.StoryTableName,
	}

	story, err := helpers.GetStoryFromRequestContext(c)
	if err != nil {
		return id, err
	}

	// Generate new id for story
	story.Id = helpers.GenerateStringUUID()
	err = table.AddStory(story)
	if err != nil {
		return id, err
	}

	id = story.Id
	log.Println("Added story: ", story.Id)
	return id, nil
}
