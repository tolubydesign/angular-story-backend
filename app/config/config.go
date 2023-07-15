package config

import (
	"github.com/joho/godotenv"
)

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
}

type Config struct {
	Configuration *DatabaseConfig
}

// Build and return the environmental configuration.
func GetConfiguration() (*Config, error) {
	// Alternative method of getting a .env file
	// gottenEnv := os.Getenv("PORT")
	// environment := os.Getenv("ENV")
	var envs map[string]string
	envs, envErr := godotenv.Read(".env")
	if envErr != nil {
		return nil, envErr
	}

	port := envs["PORT"]
	environment := envs["ENV"]

	// TODO: Get configuration settings from .env file
	// In a production environment you wouldn't hard code these values. They should be dynamically imported
	redisConfiguration := RedisConfiguration{
		User:     "",
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Database: 0,
	}

	postgreSQLConfig := PostgreSQLConfiguration{
		User:         "postgres",
		Host:         "localhost",
		Port:         5432,
		Password:     "postgres",
		DatabaseName: "postgres",
	}

	return &Config{
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
		},
	}, nil
}
