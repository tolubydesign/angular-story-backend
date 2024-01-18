package main

import (
	"context"
	"errors"
	"strings"

	// "fmt"

	// "log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/aws/aws-sdk-go-v2/config"
)

// A simple token-based authorizer example to demonstrate how to use an authorization token
// to allow or deny a request. In this example, the caller named 'user' is allowed to invoke
// a request if the client-supplied token value is 'allow'. The caller is not allowed to invoke
// the request if the token value is 'deny'. If the token value is 'unauthorized' or an empty
// string, the authorizer function returns an HTTP 401 status code. For any other token value,
// the authorizer returns an HTTP 500 status code.
// Note that token values are case-sensitive.

// https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-output.html
// https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-lambda-authorizer.html

var apiKey = os.Getenv("API_KEY")

func init() {
}

func main() {
	lambda.Start(HandleRequest)
}

// func HandleRequest(ctx context.Context, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
func HandleRequest(ctx context.Context, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.Headers["Authorization"]
	tokenSlice := strings.Split(token, " ")
	var bearerToken string
	if len(tokenSlice) > 1 {
		bearerToken = tokenSlice[len(tokenSlice)-1]
	}
	if bearerToken != apiKey {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Authentication key provide is invalid")
	}
	return generatePolicy("user", "Allow", request.MethodArn), nil

	//
	// /////...............
	//

	// var bearerToken string
	// if len(tokenSlice) > 1 {
	// 	bearerToken = tokenSlice[len(tokenSlice)-1]
	// }

	// fmt.Println("context", ctx)
	// fmt.Println("request", request)

	// // Check that user has provided the correct authentication key
	// if apiKey == "" {
	// 	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Api Key held by system")
	// }

	// if bearerToken != apiKey {
	// 	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Authentication key provide is invalid")
	// }

	// ***************
	// return generatePolicy("user", "Allow", request.MethodArn), nil
}

// *************** below is useful
// // func HandleRequest(ctx context.Context, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
// func HandleRequest(ctx context.Context, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
// 	// awsConfig, err := config.LoadDefaultConfig(context.TODO())
// 	// if err != nil {
// 	// 	log.Println("aws config failed to load.", err)

// 	// 	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
// 	// 		IsAuthorized: false,
// 	// 	}, fmt.Errorf("Invalid configuration.") // OR }, nil
// 	// }
// 	if apiKey == "" {
// 		return events.APIGatewayV2CustomAuthorizerSimpleResponse{}, errors.New("Invalid configuration")
// 	}
// 	context := map[string]interface{}{
// 		"ctx":     ctx,
// 		"request": request,
// 	}

// 	auth := request.Headers["Authorization"]
// 	if auth == "" {
// 		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
// 			IsAuthorized: false,
// 			Context:      context,
// 		}, errors.New("No authentication key found")
// 	}

// 	if apiKey != auth {
// 		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
// 			IsAuthorized: false,
// 			Context:      context,
// 		}, errors.New("Authentication key provide is invalid")
// 	}

// 	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
// 		IsAuthorized: true,
// 		Context:      context,
// 	}, nil
// }

// func HandleRequest(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
// 	token := request.AuthorizationToken
// 	tokenSlice := strings.Split(token, " ")
// 	var bearerToken string
// 	if len(tokenSlice) > 1 {
// 		bearerToken = tokenSlice[len(tokenSlice)-1]
// 	}
// 	if bearerToken != apiKey {
// 		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
// 	}

// 	return generatePolicy("user", "Allow", request.MethodArn), nil
// }

func generatePolicy(principalID string, effect string, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}
	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	// Optional output with custom properties of the String, Number or Boolean type.
	authResponse.Context = map[string]interface{}{
		"resource": resource,
		"effect":   effect,
	}

	return authResponse
}
