package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"

	// "github.com/tolubydesign/angular-story-backend/app/helpers"
	configuration "github.com/tolubydesign/angular-story-backend/app/config"
	"github.com/tolubydesign/angular-story-backend/app/models"
)

type MyEvent struct {
	Name string `json:"name"`
}

// env variables
var apiKey = os.Getenv("API_KEY")
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

// resource https://blog.omeir.dev/building-a-serverless-rest-api-with-go-aws-lambda-and-api-gateway
// Function runs before main()
func init() {
	customConfig = configuration.GenerateConfiguration(&configuration.Config{
		Configuration: &configuration.DatabaseConfig{
			Environment: "production",
			Port:        "",
			Charset:     "utf8",
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

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	db = *dynamodb.NewFromConfig(sdkConfig)
}

func main() {
	app := fiber.New()
	// app.Get("/", handlers.HealthCheck)
	// app.Get("/users", handlers.ReturnUsers)
	app.Get("/", func(c *fiber.Ctx) error {
		c.Response().StatusCode()
		c.Response().Header.Add("Content-Type", "application/json")
		return c.JSON(models.HTTPResponse{
			Code:    fiber.StatusOK,
			Message: "Lambda function is working correctly",
		})
	})

	fmt.Printf("%s is %s. years old\n", os.Getenv("NAME"), os.Getenv("AGE"))

	fiberLambda = fiberadapter.New(app)
	lambda.Start(Handler)
	// lambda.Start(HandleRequest)
	// if helpers.IsLambda() {
	// 	fiberLambda = fiberadapter.New(app)
	// 	lambda.Start(Handler)
	// } else {
	// 	app.Listen(":3000")
	// }
}

func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}
	message := fmt.Sprintf("Hello %s!", event.Name) + fmt.Sprintf("Successful Request Tolu!")
	return &message, nil
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return fiberLambda.ProxyWithContext(ctx, request)
}
