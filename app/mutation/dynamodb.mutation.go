package mutation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	models "github.com/tolubydesign/angular-story-backend/app/models"
)

// Reading:
// [Function declaration syntax: things in parenthesis before function name](https://stackoverflow.com/questions/34031801/function-declaration-syntax-things-in-parenthesis-before-function-name)
// [DynamoDB examples using SDK for Go V2](https://docs.aws.amazon.com/code-library/latest/ug/go_2_dynamodb_code_examples.html)

// TableBasics encapsulates the Amazon DynamoDB service actions used in the examples.
// It contains a DynamoDB service client that is used to act on the specified table.
type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type CreateTableStruct struct {
	AttributeDefinition []types.AttributeDefinition
	KeySchemaElement    []types.KeySchemaElement
}

// Determines whether a DynamoDB table exists.
func (client TableBasics) TableExists() (bool, error) {
	exists := true
	_, err := client.DynamoDbClient.DescribeTable(
		context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(client.TableName)},
	)

	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			log.Printf("Table %v does not exist.\n", client.TableName)
			err = nil
		} else {
			log.Printf("Couldn't determine existence of table %v. Here's why: %v\n", client.TableName, err)
		}

		exists = false
	}

	return exists, err
}

// CreateMovieTable creates a DynamoDB table with a composite primary key defined as
// a string sort key named `title`, and a numeric partition key named `year`.
// This function uses NewTableExistsWaiter to wait for the table to be created by
// DynamoDB before it returns.
func (basics TableBasics) CreateMovieTable() (*types.TableDescription, error) {
	var tableDesc *types.TableDescription
	table, err := basics.DynamoDbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("year"),
			AttributeType: types.ScalarAttributeTypeN,
		}, {
			AttributeName: aws.String("title"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("year"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("title"),
			KeyType:       types.KeyTypeRange,
		}},
		TableName: aws.String(basics.TableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	if err != nil {
		log.Printf("\nERROR: Couldn't create table %v. ERROR - %v\n", basics.TableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(basics.DynamoDbClient)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(basics.TableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}
	return tableDesc, err
}

/*
Create dynamodb table if it does not already exist

Provide name of table within the client.

Create a DynamoDB table with a composite primary key defined as
a string sort key named `title`, and a numeric partition key named `year`.
This function uses NewTableExistsWaiter to wait for the table to be created by
DynamoDB before it returns.
*/
func (client TableBasics) CreateDynamoDBTable(tableStruct CreateTableStruct) (*types.TableDescription, error) {
	var tableDesc *types.TableDescription
	table, err := client.DynamoDbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{

		AttributeDefinitions: tableStruct.AttributeDefinition,
		KeySchema:            tableStruct.KeySchemaElement,

		TableName: aws.String(client.TableName),

		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})

	if err != nil {
		// Log information
		log.Printf("\nERROR: Create DynamoDB Table, name:_ %v. message:_ %v\n", client.TableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(client.DynamoDbClient)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(client.TableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Reasoning: %v\n", err)
		}
		tableDesc = table.TableDescription
	}

	return tableDesc, err
}

/*
Removes an existing table from the DynamoDB.
{...}
Returns possible error, if table does not exist
*/
func DeleteDynamoDBTable(client *dynamodb.Client) error {
	var err error
	// {...}
	return err
}

// TODO: add multiple
// Add a story the DynamoDB table.
func (basics TableBasics) AddStory(story models.DynamoStoryDatabaseStruct) error {
	fmt.Println("Adding story to database.")

	item, err := attributevalue.MarshalMap(story)
	if err != nil {
		return err
	}
	_, err = basics.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: item,
	})

	if err != nil {
		log.Printf("Couldn't add item to table. Reasoning: %v\n", err)
	}

	return err
}
