package database

import (
	"database/sql"
	"errors"
	"fmt"
)

var postgreSQLDatabaseSingleton *sql.DB

func ConnectToPostgreSQLDatabase() (*sql.DB, error) {
	var err error
	// TODO: production and development relevant
	const (
		postgresHost     = "localhost"
		postgresPort     = 5432
		postgresUser     = "postgres"
		postgresPassword = "postgres"
		postgresDBname   = "postgres"
	)

	connection := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDBname)

	// Connect to database
	postgreSQLDatabaseSingleton, err = sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	return postgreSQLDatabaseSingleton, nil
}

func GetPostgreSQLDatabaseSingleton() (*sql.DB, error) {
	if postgreSQLDatabaseSingleton == nil {
		return nil, errors.New("PostgreSQL Database singleton is null")
	}

	return postgreSQLDatabaseSingleton, nil
}
