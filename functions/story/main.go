package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"

	configuration "github.com/tolubydesign/angular-story-backend/app/config"
	dynamodbrequest "github.com/tolubydesign/angular-story-backend/app/controller/dynamodb-request"
	"github.com/tolubydesign/angular-story-backend/app/models"
	"github.com/tolubydesign/angular-story-backend/app/utils"
)

type MyEvent struct {
	Name string `json:"name"`
}

// env variables
var env = os.Getenv("ENV")
var JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
var dynamodbStoryTableName = os.Getenv("DYNAMODB_STORY_TABLE_NAME")
var dynamodbUserTableName = os.Getenv("DYNAMODB_USER_TABLE_NAME")
var accessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
var secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
var region = os.Getenv("AWS_REGION")
var accountID = os.Getenv("AWS_ACCOUNT_ID")
var sessionToken = os.Getenv("AWS_SESSION_TOKEN")

var fiberLambda *fiberadapter.FiberLambda
var customConfig *configuration.Config
var dynamodbClient *dynamodb.Client

// Function related variables

// resource https://blog.omeir.dev/building-a-serverless-rest-api-with-go-aws-lambda-and-api-gateway
func init() {
	customConfig = configuration.GenerateConfiguration(&configuration.Config{
		Configuration: &configuration.DatabaseConfig{
			Environment: env,
			// Port:        "",
			Charset: "utf8",
			Redis: configuration.RedisConfiguration{
				User:     "",
				Host:     "",
				Port:     0000,
				Password: "",
				Database: 0000,
			},
			Dynamodb: configuration.DynamodbConfiguration{
				StoryTableName: dynamodbStoryTableName,
				UserTableName:  dynamodbUserTableName,
			},
			JWTSecretKey: []byte(JWTSecretKey),
			AWS: configuration.AWSConfiguration{
				AccessKeyID:     accessKeyID,
				SecretAccessKey: secretAccessKey,
				SessionToken:    sessionToken,
				Region:          region,
				AccountID:       accountID,
			},
		},
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		log.Printf("error: %s", err.Error())
		return
		// log.Fatal("unable to load SDK config, %v", err)
	}

	dynamodbClient = dynamodb.NewFromConfig(cfg)
}

func main() {
	if dynamodbClient == nil {
		log.Fatal("Dynamodb Client is not connected.")
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Response().StatusCode()
		c.Response().Header.Add("Content-Type", "application/json")
		return c.JSON(models.HTTPResponse{
			Code:    fiber.StatusOK,
			Message: "lambda function accessible",
		})
	})

	app.Get(utils.Endpoints.Get.Story, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.GetStoryByIdRequest(ctx, dynamodbClient, customConfig)
	})

	app.Get(utils.Endpoints.Get.AllStories, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.ListAllStoriesRequest(ctx, dynamodbClient, customConfig)
	})

	app.Post(utils.Endpoints.Post.Story, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.AddStoryRequest(ctx, dynamodbClient, customConfig)
	})

	app.Put(utils.Endpoints.Put.Story, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.UpdateDynamodbStoryRequest(ctx, dynamodbClient, customConfig)
	})

	app.Delete(utils.Endpoints.Delete.Story, func(ctx *fiber.Ctx) error {
		return dynamodbrequest.DeleteDynamodbStoryRequest(ctx, dynamodbClient, customConfig)
	})

	fiberLambda = fiberadapter.New(app)
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return fiberLambda.ProxyWithContext(ctx, request)
}
