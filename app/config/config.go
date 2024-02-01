package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AWSConfiguration struct {
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	Region          string `json:"region"`
	AccountID       string `json:"accountID"`
}

type DynamodbConfiguration struct {
	StoryTableName string `json:"storyTableName"`
	UserTableName  string `json:"userTableName"`
}

type RedisConfiguration struct {
	User     string `json:"user"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

type DatabaseConfig struct {
	Environment  string                `json:"environment"`
	Port         string                `json:"port"`
	Charset      string                `json:"charset"`
	Redis        RedisConfiguration    `json:"redis"`
	Dynamodb     DynamodbConfiguration `json:"dynamodb"`
	JWTSecretKey []byte                `json:"jwtSecretKey"`
	AWS          AWSConfiguration      `json:"aws"`
}

type Config struct {
	Configuration *DatabaseConfig
}

var configurationSingleton *Config

/*
Build and return the environmental configuration.

Returns Configuration or error, if issues occur.
*/
func BuildConfiguration() (*Config, error) {
	envArg := os.Args[1]
	environmentPath := fmt.Sprintf(".env.%s", envArg)
	log.Println("environment: ", envArg, " | .env file: ", environmentPath)

	// Deny processing if environment argument isn't what we want
	if (envArg == "development") || (envArg == "production") {
		// Note: Alternative method of getting a .env file
		// gottenEnv := os.Getenv("PORT")
		// environment := os.Getenv("ENV")
		var envs map[string]string
		envs, err := godotenv.Read(environmentPath)
		if err != nil {
			return nil, err
		}

		port := envs["PORT"]
		environment := envs["ENV"]
		secret := envs["JWT_SECRET_KEY"]

		aws := AWSConfiguration{
			AccessKeyID:     envs["AWS_ACCESS_KEY_ID"],
			SecretAccessKey: envs["AWS_SECRET_ACCESS_KEY"],
			SessionToken:    "dummy",
			Region:          envs["AWS_REGION"],
			AccountID:       envs["AWS_ACCOUNT_ID"],
		}

		// TODO: Get configuration settings from .env file
		dynamodbConfiguration := DynamodbConfiguration{
			StoryTableName: envs["DYNAMODB_STORY_TABLE_NAME"],
			UserTableName:  envs["DYNAMODB_USER_TABLE_NAME"],
		}

		// TODO: Get configuration from .env file
		redisConfiguration := RedisConfiguration{
			User:     "",
			Host:     "localhost",
			Port:     6379,
			Password: "",
			Database: 0,
		}

		configurationSingleton = &Config{
			Configuration: &DatabaseConfig{
				Environment:  environment,
				Port:         port,
				Charset:      "utf8",
				Redis:        redisConfiguration,
				Dynamodb:     dynamodbConfiguration,
				JWTSecretKey: []byte(secret),
				AWS:          aws,
			},
		}

		return configurationSingleton, nil
	}

	return nil, errors.New("Incorrect environment variables provided. Only 'production' or 'development' can be used")
}

func GetConfiguration() (*Config, error) {
	if configurationSingleton == nil {
		// Build configuration
		build, e := BuildConfiguration()
		if e != nil {
			return nil, errors.New("Project Configuration is undefined")
		}

		return build, nil
	}

	return configurationSingleton, nil
}

// Generate a custom configuration object based on variables passed.
// This function does not effect the configuration singleton
func GenerateConfiguration(config *Config) *Config {
	var generatedConfiguration *Config
	port := config.Configuration.Port
	environment := config.Configuration.Environment
	secret := config.Configuration.JWTSecretKey
	aws := AWSConfiguration{
		AccessKeyID:     config.Configuration.AWS.AccessKeyID,
		SecretAccessKey: config.Configuration.AWS.SecretAccessKey,
		SessionToken:    config.Configuration.AWS.SessionToken,
		Region:          config.Configuration.AWS.Region,
		AccountID:       config.Configuration.AWS.AccountID,
	}
	dynamodbConfiguration := DynamodbConfiguration{
		StoryTableName: config.Configuration.Dynamodb.StoryTableName,
		UserTableName:  config.Configuration.Dynamodb.UserTableName,
	}
	redisConfiguration := RedisConfiguration{
		User:     config.Configuration.Redis.User,
		Host:     config.Configuration.Redis.Host,
		Port:     config.Configuration.Redis.Port,
		Password: config.Configuration.Redis.Password,
		Database: config.Configuration.Redis.Database,
	}
	generatedConfiguration = &Config{
		Configuration: &DatabaseConfig{
			Environment:  environment,
			Port:         port,
			Charset:      "utf8",
			Redis:        redisConfiguration,
			Dynamodb:     dynamodbConfiguration,
			JWTSecretKey: []byte(secret),
			AWS:          aws,
		},
	}

	return generatedConfiguration
}
