package dynamodbrequest

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/tolubydesign/angular-story-backend/app/config"
	"github.com/tolubydesign/angular-story-backend/app/models"
	"github.com/tolubydesign/angular-story-backend/app/query"
	"github.com/tolubydesign/angular-story-backend/app/utils"
)

/*
Check if user credentials is a valid user. Find user via email and name
*/
func UserVerification(client *dynamodb.Client, user models.DatabaseUserStruct) (string, error) {
	var token string
	var err error

	configuration, err := config.GetConfiguration()
	if err != nil {
		return token, err
	}

	table := query.TableBasics{
		DynamoDbClient: client,
		TableName:      configuration.Configuration.Dynamodb.UserTableName,
	}

	u, err := table.GetUserByEmail(user.Email)
	if err != nil {
		return token, err
	}

	// Validate that password provided matches known (hashed) password, in database
	// If user is found. Check if password matches database password
	if u.Email != user.Email || u.Password != user.Password {
		return token, errors.New("User information doesn't match known user information.")
	}
	
	// Set Token
	token, err = utils.BuildUserJWTToken(u)
	if err != nil {
		return token, err
	}

	return token, err
}
