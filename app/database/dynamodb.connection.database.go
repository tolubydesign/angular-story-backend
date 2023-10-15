package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	envConfiguration "github.com/tolubydesign/angular-story-backend/app/config"
	"github.com/tolubydesign/angular-story-backend/app/helpers"
	"github.com/tolubydesign/angular-story-backend/app/mutation"
	"github.com/tolubydesign/angular-story-backend/app/query"
)

var dynamoSingleton *dynamodb.Client

/*
Connecting to dynamo db through terminal https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.UsageNotes.html
*/
func CreateDynamoClient() (*dynamodb.Client, error) {
	configuration, err := envConfiguration.GetConfiguration()
	if err != nil {
		return nil, err
	}

	env := configuration.Configuration.Environment

	var dynamoConfigEndpoint config.LoadOptionsFunc
	var dynamoConfigCredentialProvider config.LoadOptionsFunc
	var dynamoConfigWithRegion config.LoadOptionsFunc

	if env == "development" {
		// Development endpoint
		dynamoConfigWithRegion = config.WithRegion("us-east-2")

		dynamoConfigEndpoint = config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			}),
		)

		dynamoConfigCredentialProvider = config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			// TODO: use env configuration
			Value: aws.Credentials{
				AccessKeyID: "DUMMYIDEXAMPLE", SecretAccessKey: "DUMMYEXAMPLEKEY", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		})
	} else {
		// Production endpoint
		dynamoConfigWithRegion = config.WithRegion("us-east-2")

		dynamoConfigEndpoint = config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			}),
		)

		dynamoConfigCredentialProvider = config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "DUMMYIDEXAMPLE", SecretAccessKey: "DUMMYEXAMPLEKEY", SessionToken: "dummy",
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

	// setup database with dummy data
	if env == "development" {
		tableFound, err := StoryTableExists(dynamoSingleton)
		if err != nil {
			return dynamoSingleton, err
		}

		if tableFound == false {
			fmt.Println("Table 'Story' table was not found. Attempting to create a new table")
			SetupStoryDatabase(dynamoSingleton)
		}
	}

	return dynamoSingleton, nil
}

func GetDynamoSingleton() (*dynamodb.Client, error) {
	if dynamoSingleton == nil {
		return nil, errors.New("Dynamo Client was not accessible.")
	}

	return dynamoSingleton, nil
}

/*
Check that the required Story table exists in the database.
*/
func StoryTableExists(client *dynamodb.Client) (bool, error) {
	var exists = false
	tableName := "Story"
	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}

	tables, err := table.ListDynamodbTables()

	if err != nil {
		return false, err
	}

	// Check for story table in array of table names
	for index, a := range tables {
		fmt.Println("looping through tables ", index, a)
		if a == "Story" {
			exists = true
		}
	}

	return exists, err
}

/*
For development purposes only
Function sets up database with relevant tables and data for consumption but end user.
Function will input dummy data.

Critical research
[Dynamodb- Adding non key attributes](https://stackoverflow.com/questions/38151687/dynamodb-adding-non-key-attributes)
*/
func SetupStoryDatabase(client *dynamodb.Client) error {
	var err error
	if client == nil {
		return errors.New("Database is unreachable")
	}

	// Reading
	// [Pointer receivers](https://go.dev/tour/methods/4)
	tableName := "Story"
	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      *aws.String(tableName),
	}

	// Table structure
	fmt.Println("Setting up Story Database structure")
	tableDescription, err := table.CreateDynamoDBTable(mutation.CreateTableStruct{
		// Add attribute definition
		AttributeDefinition: []types.AttributeDefinition{{
			AttributeName: aws.String("storyId"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("title"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		// Add key schema
		KeySchemaElement: []types.KeySchemaElement{{
			AttributeName: aws.String("storyId"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("title"),
			KeyType:       types.KeyTypeRange,
		}},
	})

	if err != nil {
		return err
	}

	fmt.Printf("\n Information about newly created dynamodb table, %v", tableDescription)
	err = helpers.PopulateStoryDatabase(table)
	if err != nil {
		return err
	}

	return nil
}
