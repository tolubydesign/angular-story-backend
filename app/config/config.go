package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AWSConfiguration struct {
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
}

type DynamodbConfiguration struct {
	Aws            AWSConfiguration `json:"aws"`
	StoryTableName string           `json:"storyTableName"`
	UserTableName  string           `json:"userTableName"`
}

type RedisConfiguration struct {
	User     string `json:"user"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

type PostgreSQLConfiguration struct {
	User         string `json:"user"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Password     string `json:"password"`
	DatabaseName string `json:"databaseName"`
}

type DatabaseConfig struct {
	Environment string                  `json:"environment"`
	Port        string                  `json:"port"`
	Dialect     string                  `json:"dialect"`
	Username    string                  `json:"username"`
	Password    string                  `json:"password"`
	Name        string                  `json:"name"`
	Charset     string                  `json:"charset"`
	Redis       RedisConfiguration      `json:"redis"`
	Postgres    PostgreSQLConfiguration `json:"postgres"`
	Dynamodb    DynamodbConfiguration   `json:"dynamodb"`
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
	pgArgEnvironment := os.Args[1]
	environmentPath := fmt.Sprintf(".env.%s", pgArgEnvironment)

	// Deny processing if environment argument isn't what we want
	if (pgArgEnvironment == "development") || (pgArgEnvironment == "production") {
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

		// TODO: Get configuration settings from .env file
		dynamodbConfiguration := DynamodbConfiguration{
			Aws: AWSConfiguration{
				AccessKeyID:     envs["AWS_ACCESS_KEY_ID"],
				SecretAccessKey: envs["AWS_SECRET_ACCESS_KEY"],
				SessionToken:    "dummy",
			},
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

		// TODO: Get configuration from .env file
		postgreSQLConfig := PostgreSQLConfiguration{
			User:         "postgres",
			Host:         "localhost",
			Port:         5432,
			Password:     "postgres",
			DatabaseName: "postgres",
		}

		configurationSingleton = &Config{
			Configuration: &DatabaseConfig{
				Environment: environment,
				Port:        port,
				Dialect:     "mysql",
				Username:    "root",
				Password:    "",
				Name:        "testinger",
				Charset:     "utf8",
				Redis:       redisConfiguration,
				Postgres:    postgreSQLConfig,
				Dynamodb:    dynamodbConfiguration,
			},
		}

		return configurationSingleton, nil
	}

	return nil, errors.New("Incorrect environment variables provided. Only 'production' or 'development' can be used")
}

func GetConfiguration() (*Config, error) {
	if configurationSingleton == nil {
		return nil, errors.New("Project Configuration is undefined")
	}

	return configurationSingleton, nil
}
