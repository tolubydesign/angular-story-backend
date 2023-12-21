package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/tolubydesign/angular-story-backend/app/helpers"
	"github.com/tolubydesign/angular-story-backend/app/models"
)

var weekday string
var fiberLambda *fiberadapter.FiberLambda

// resource https://blog.omeir.dev/building-a-serverless-rest-api-with-go-aws-lambda-and-api-gateway
// Function runs before main()
func init() {
	weekday = time.Now().Weekday().String()
}

func main() {
	fmt.Printf("Today is %s", weekday)

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
		// return c.SendString("Lambda function is working correctly")
	})

	if helpers.IsLambda() {
		fiberLambda = fiberadapter.New(app)
		lambda.Start(Handler)
	} else {
		app.Listen(":3000")
	}
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return fiberLambda.ProxyWithContext(ctx, request)
}
