package query

import (
	"context"
	"errors"
	"log"

	"github.com/tolubydesign/angular-story-backend/app/models"
	"github.com/tolubydesign/angular-story-backend/app/mutation"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Code
// https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/dynamodb/actions/table_basics.go#L272
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Scan.html
// https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/go/example_code/dynamodb/DynamoDBCreateItem.go
// https://aws.github.io/aws-sdk-go-v2/docs/getting-started/
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/getting-started-step-5.html
// https://davidagood.com/dynamodb-local-go/

// setting type
type TableBasics mutation.TableBasics

// List the DynamoDB table names for the current account.
func (basics TableBasics) ListDynamodbTables() ([]string, error) {
	var tableNames []string
	tables, err := basics.DynamoDbClient.ListTables(
		context.TODO(),
		&dynamodb.ListTablesInput{},
	)

	if err != nil {
		log.Printf("Couldn't list tables. Here's why: %v\n", err)
	} else {
		tableNames = tables.TableNames
	}
	return tableNames, err
}

// Scan gets all movies in the DynamoDB table that were released in a range of years
// and projects them to return a reduced set of fields.
// The function uses the `expression` package to build the filter and projection
// expressions.
func (basics TableBasics) ScanByExpression(startYear int, endYear int) ([]models.DynamoStoryResponseStruct, error) {
	var stories []models.DynamoStoryResponseStruct
	var err error
	var response *dynamodb.ScanOutput

	filtEx := expression.Name("year").Between(expression.Value(startYear), expression.Value(endYear))
	projEx := expression.NamesList(
		expression.Name("year"), expression.Name("title"), expression.Name("info.rating"),
	)

	expr, err := expression.NewBuilder().WithFilter(filtEx).WithProjection(projEx).Build()

	if err != nil {
		log.Printf("Couldn't build expressions for scan. Here's why: %v\n", err)
	} else {
		response, err = basics.DynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName:                 aws.String(basics.TableName),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection(),
		})

		if err != nil {
			log.Printf("Couldn't scan for movies released between %v and %v. Error: %v\n",
				startYear, endYear, err)
		} else {
			err = attributevalue.UnmarshalListOfMaps(response.Items, &stories)
			if err != nil {
				log.Printf("Couldn't unmarshal query response. Error: %v\n", err)
			}
		}
	}

	return stories, err
}

/*
Scan the entire Story table.
Return all items in table.
*/
func (basics TableBasics) FullTableScan() ([]models.DynamoStoryResponseStruct, error) {
	var stories []models.DynamoStoryResponseStruct
	var err error

	response, err := basics.DynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:              aws.String(basics.TableName),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	})

	if err != nil {
		// log.Printf("\nCouldn't scan for stories. Error: %v\n", err)
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &stories)
	if err != nil {
		// log.Printf("Couldn't unmarshal query response. Error: %v\n", err)
		return nil, err
	}

	if response.LastEvaluatedKey != nil {
		log.Println("All items have not been scanned")
		return nil, errors.New("All items have not been scanned")
	}

	return stories, nil
}

func (basics TableBasics) GetStoryById(id string) (*models.DynamoStoryResponseStruct, error) {
	var stories []*models.DynamoStoryResponseStruct
	var err error
	var response *dynamodb.ScanOutput

	nameBuilder := expression.Name("id").Equal(expression.Value(id))
	// projectBuilder := expression.NamesList(
	// 	expression.Name("year"), expression.Name("title"), expression.Name("info.rating"),
	// )

	expr, err := expression.NewBuilder().WithFilter(nameBuilder).Build()

	if err != nil {
		log.Printf("Couldn't build expressions for scan. Here's why: %v\n", err)
		return nil, err
	} else {
		response, err = basics.DynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName:                 aws.String(basics.TableName),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			// ProjectionExpression:      expr.Projection(),
		})

		if err != nil {
			log.Printf("\nCouldn't scan for story with id %v Error: %v\n", id, err)
			return nil, err
		} else {
			err = attributevalue.UnmarshalListOfMaps(response.Items, &stories)
			if err != nil {
				log.Printf("Couldn't unmarshal query response. Error: %v\n", err)
				return nil, err
			}
		}
	}

	if len(stories) == 0 {
		return nil, errors.New("No user found.")
	}

	return stories[0], nil
}

// Scan database for users
func (basics TableBasics) FullUserTableScan() ([]models.DatabaseUserStruct, error) {
	var users []models.DatabaseUserStruct
	var err error

	response, err := basics.DynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:              aws.String(basics.TableName),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	})

	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &users)
	if err != nil {
		return nil, err
	}

	if response.LastEvaluatedKey != nil {
		log.Println("All items have not been scanned")
		return nil, errors.New("All items have not been scanned")
	}

	if len(users) == 0 {
		return nil, errors.New("No data found.")
	}

	return users, nil
}

func (basics TableBasics) GetUserByEmail(email string) (*models.DatabaseUserStruct, error) {
	// TODO: security log
	var users []*models.DatabaseUserStruct
	var err error
	var response *dynamodb.ScanOutput
	nameBuilder := expression.Name("email").Equal(expression.Value(email))
	expr, err := expression.NewBuilder().WithFilter(nameBuilder).Build()

	if err != nil {
		log.Printf("Couldn't build expressions for scan. Here's why: %v\n", err)
		return nil, err
	}

	response, err = basics.DynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:                 aws.String(basics.TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})

	if err != nil {
		log.Printf("\nCouldn't scan for story with id %v Error: %v\n", email, err)
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &users)
	if err != nil {
		log.Printf("Couldn't unmarshal query response. Error: %v\n", err)
		return nil, err
	}

	if len(users) == 0 {
		// TODO: create more ambiguous error response
		return nil, errors.New("No data found.")
	}

	return users[0], nil
}
