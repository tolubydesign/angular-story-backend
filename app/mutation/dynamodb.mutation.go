package mutation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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

// type DataBaseDynamoStoryStruct models.DynamoStoryDatabaseStruct

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

// Add a story the DynamoDB table.
func (basics TableBasics) AddStory(story models.DynamoStoryDatabaseStruct) error {
	log.Println("Database interaction. Adding story to database")
	item, err := attributevalue.MarshalMap(story)
	if err != nil {
		return err
	}
	output, err := basics.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: item,
	})

	log.Println("DynamoDB Client Output: ", output)
	if err != nil {
		log.Println("Couldn't add item to table. Reasoning: ", err.Error())
		return err
	}

	return err
}

func (basics TableBasics) AddUser(user models.DatabaseUserStruct) error {
	log.Println("Add user with id:", user.Id)

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = basics.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: item,
	})

	if err != nil {
		log.Printf("Couldn't add user to table. Reasoning: %v\n", err)
		return err
	}

	return err
}

/*
Update existing story in the dynamodb database with new content provided.
This function uses the `expression` package to build the update expression.

resource - https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GettingStarted.UpdateItem.html \
resource - https://dave.dev/blog/2022/07/06-08-2022-ddbtools/ \
resource - https://yuminlee2.medium.com/golang-access-struct-fields-ae320fb74d17 \
resource - https://dynobase.dev/dynamodb-golang-query-examples/#update-item \
*/
func (basics TableBasics) UpdateDynamoDBTable(story models.DynamoStoryDatabaseStruct) (models.DynamoStoryDatabaseStruct, error) {
	log.Println("Updating existing story, using story content provided.")
	var err error
	var response *dynamodb.UpdateItemOutput
	var attributeMap models.DynamoStoryDatabaseStruct

	pointer := &models.DynamoStoryDatabaseStruct{
		Id:          story.Id,
		Creator:     story.Creator,
		Title:       story.Title,
		Description: story.Description,
		Content:     story.Content,
	}

	// Note: get `types.AttributeValue` value for value provided
	// t, err := attributevalue.Marshal(story.Title)
	// if err != nil {
	// 	return attributeMap, err
	// }

	update := expression.Set(expression.Name("description"), expression.Value((*pointer).Description))
	update.Set(expression.Name("title"), expression.Value((*pointer).Title))
	update.Set(expression.Name("content"), expression.Value((*pointer).Content))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		return attributeMap, err
	} else {
		response, err = basics.DynamoDbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName: aws.String(basics.TableName),
			Key: map[string]types.AttributeValue{
				"id":      &types.AttributeValueMemberS{Value: (*pointer).Id},
				"creator": &types.AttributeValueMemberS{Value: (*pointer).Creator},
			},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,

			// Note: alternative method of making dynamodb request
			// UpdateExpression:         aws.String("set description = :description, content = :content"),
			// ExpressionAttributeValues: map[string]types.AttributeValue{
			// 	":title":       &types.AttributeValueMemberS{Value: (*pointer).Title},
			// 	":description": &types.AttributeValueMemberS{Value: (*pointer).Description},
			// 	":content":     storyContent,
			// },
		})

		consoleResponse := fmt.Sprintf("%v", response)
		log.Println("Update story return response: ", consoleResponse)

		if err != nil {
			log.Printf("Couldn't update story %v. Here's why: %v\n", story.Title, err)
			return attributeMap, err
		} else {

			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				log.Printf("Couldn't unmarshall update response. Here's why: %v\n", err)
				return attributeMap, err
			}

			return attributeMap, err
		}
	}
}

// DeleteMovie removes a movie from the DynamoDB table.
func (basics TableBasics) DeleteStory(story models.DynamoStoryDatabaseStruct) error {
	_, err := basics.DynamoDbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(basics.TableName),
		Key: map[string]types.AttributeValue{
			"id":      &types.AttributeValueMemberS{Value: story.Id},
			"creator": &types.AttributeValueMemberS{Value: story.Creator},
		},
	})

	// [Knowledge] how to show print data in human readable json
	// d, err := json.Marshal(output)
	// log.Println("output: ", string(d))

	if err != nil {
		// log.Println("Couldn't delete id: ", story.Id, " creator: ", story.Creator)
		log.Println("Couldn't delete id: ", story.Id, " reasoning: ", err.Error())
		return err
	}

	return nil
}

// GetKey returns the composite primary key of the movie in a format that can be
// sent to DynamoDB.
// resource - https://github.com/awsdocs/aws-doc-sdk-examples/blob/1c12c397d9bf042f81194ce0621fb443d4712317/gov2/dynamodb/actions/movie.go
func GetKey(story models.DynamoStoryDatabaseStruct) (map[string]types.AttributeValue, error) {
	title, err := attributevalue.Marshal(story.Title)
	if err != nil {
		return map[string]types.AttributeValue{"title": title}, err
	}

	return map[string]types.AttributeValue{"title": title}, err
}
