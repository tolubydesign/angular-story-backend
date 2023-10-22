package query

import (
	"context"
	"errors"
	"log"

	"github.com/tolubydesign/angular-story-backend/app/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tolubydesign/angular-story-backend/app/mutation"
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
func (basics TableBasics) ScanByExpression(startYear int, endYear int) ([]models.DynamoStoryStruct, error) {
	var stories []models.DynamoStoryStruct
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

// TODO: add description
// {...}
func (basics TableBasics) FullTableScan() ([]models.DynamoStoryStruct, error) {
	var stories []models.DynamoStoryStruct
	var err error

	response, err := basics.DynamoDbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:              aws.String(basics.TableName),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	})

	// if err != nil {
	// 	log.Fatal("Scan failed", err)
	// }

	if err != nil {
		// log.Printf("\nCouldn't scan for stories. Error: %v\n", err)
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &stories)
	if err != nil {
		// log.Printf("Couldn't unmarshal query response. Error: %v\n", err)
		return nil, err
	}

	// for _, i := range response.Items {
	// 	var u models.DynamoStoryStruct
	// 	err := attributevalue.UnmarshalMap(i, &u)
	// 	if err != nil {
	// 		log.Fatal("unmarshal failed", err)
	// 	}
	// }

	if response.LastEvaluatedKey != nil {
		log.Println("all items have not been scanned")
		return nil, errors.New("All items have not been scanned")
	}

	// log.Println("scanned", response.ScannedCount, "items in", time.Since(startTime).Seconds(), "seconds")
	// log.Println("consumed capacity", *response.ConsumedCapacity.CapacityUnits)

	return stories, nil
}
