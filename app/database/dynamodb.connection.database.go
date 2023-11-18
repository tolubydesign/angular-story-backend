package database

import (
	"context"
	"errors"
	"fmt"
	"log"

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
		// Check that story table is in dynamo db
		tableFound, err := StoryTableExists(dynamoSingleton)
		if err != nil {
			return dynamoSingleton, err
		}

		if tableFound == false {
			fmt.Println("Table 'Story' table was not found. Attempting to create a new Story table.")
			SetupStoryDatabase(dynamoSingleton)
		}

		// Check that user table is in the dynamo db
		userTableFound, err := UserTableExists(dynamoSingleton)
		if err != nil {
			return dynamoSingleton, err
		}

		if userTableFound == false {
			log.Println("Table 'User' table was not found. Attempting to create a new User table.")
			SetupUserData(dynamoSingleton)
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
	configuration, err := envConfiguration.GetConfiguration()
	name := configuration.Configuration.Dynamodb.StoryTableName
	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      name,
	}

	tables, err := table.ListDynamodbTables()

	if err != nil {
		return false, err
	}

	// Check for story table in array of table names
	for _, a := range tables {
		if a == "Story" {
			exists = true
		}
	}

	if exists {
		log.Println("'Story' dynamo table exists")
	}

	return exists, err
}

// Check that User Dynamo db Table exists.
func UserTableExists(client *dynamodb.Client) (bool, error) {
	var exists = false
	configuration, err := envConfiguration.GetConfiguration()
	name := configuration.Configuration.Dynamodb.UserTableName
	if err != nil {
		return exists, err
	}

	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      name,
	}

	tables, err := table.ListDynamodbTables()
	if err != nil {
		return exists, err
	}

	// Check for story table in array of table names
	for _, a := range tables {
		if a == name {
			exists = true
		}
	}

	if exists {
		log.Println("'User' dynamo table exists")
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
	configuration, err := envConfiguration.GetConfiguration()
	if client == nil {
		return errors.New("Database is unreachable")
	}

	// Reading:
	// [Pointer receivers](https://go.dev/tour/methods/4)
	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      *aws.String(configuration.Configuration.Dynamodb.StoryTableName),
	}

	// Table structure
	fmt.Println("Setting up Story Database structure")
	description, err := table.CreateDynamoDBTable(mutation.CreateTableStruct{
		// Add attribute definition
		AttributeDefinition: []types.AttributeDefinition{{
			AttributeName: aws.String("id"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("creator"),
			AttributeType: types.ScalarAttributeTypeS,
		}},

		// Add key schema
		KeySchemaElement: []types.KeySchemaElement{{
			AttributeName: aws.String("id"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("creator"),
			KeyType:       types.KeyTypeRange,
		}},
	})

	if err != nil {
		return err
	}

	fmt.Printf("\nInformation about newly created dynamodb table, %v", description)
	err = helpers.PopulateStoryDatabase(table)
	if err != nil {
		return err
	}

	return err
}

// Create User table structure, in Dynamo database.
func SetupUserData(client *dynamodb.Client) error {
	var err error
	configuration, err := envConfiguration.GetConfiguration()
	if client == nil {
		return errors.New("Database is unreachable")
	}

	table := mutation.TableBasics{
		DynamoDbClient: client,
		TableName:      *aws.String(configuration.Configuration.Dynamodb.UserTableName),
	}

	// Setup User table structure
	log.Println("Setting up User Database structure")
	_, err = table.CreateDynamoDBTable(mutation.CreateTableStruct{
		// Add attribute definition
		AttributeDefinition: []types.AttributeDefinition{{
			AttributeName: aws.String("id"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("email"),
			AttributeType: types.ScalarAttributeTypeS,
		}},

		// Add key schema
		KeySchemaElement: []types.KeySchemaElement{{
			AttributeName: aws.String("id"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("email"),
			KeyType:       types.KeyTypeRange,
		}},
	})
	if err != nil {
		return err
	}

	err = helpers.PopulateUserDatabase(table)
	if err != nil {
		return err
	}

	log.Println("User database structure was generated")
	return nil
}
