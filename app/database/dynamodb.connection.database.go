package database

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	envConfiguration "github.com/tolubydesign/angular-story-backend/app/config"
)

var dynamoSingleton *dynamodb.Client

/*
Connecting to dynamo db through terminal https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.UsageNotes.html
*/
func CreateDynamoClient() (*dynamodb.Client, error) {
	configuration, err := envConfiguration.GetProjectConfiguration()
	if err != nil {
		return nil, err
	}

	environmentConfig := configuration.Configuration.Environment

	var dynamoConfigEndpoint config.LoadOptionsFunc
	var dynamoConfigCredentialProvider config.LoadOptionsFunc
	var dynamoConfigWithRegion config.LoadOptionsFunc

	if environmentConfig == "development" {
		// Development endpoint
		dynamoConfigWithRegion = config.WithRegion("us-east-1")

		dynamoConfigEndpoint = config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			}),
		)

		dynamoConfigCredentialProvider = config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		})
	} else {
		// Production endpoint
		dynamoConfigWithRegion = config.WithRegion("us-east-1")

		dynamoConfigEndpoint = config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			}),
		)

		dynamoConfigCredentialProvider = config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		})
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		dynamoConfigWithRegion,
		dynamoConfigEndpoint,
		dynamoConfigCredentialProvider,
	)

	if err != nil {
		return nil, err
	}

	dynamoSingleton = dynamodb.NewFromConfig(cfg)
	return dynamoSingleton, nil
}

func GetDynamoSingleton() (*dynamodb.Client, error) {
	if dynamoSingleton == nil {
		return nil, errors.New("Dynamo Client was not accessible.")
	}

	return dynamoSingleton, nil
}
