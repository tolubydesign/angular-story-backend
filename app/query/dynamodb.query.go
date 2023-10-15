package query

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/tolubydesign/angular-story-backend/app/mutation"
)

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
