package models

import "database/sql"

type sqlDatabase struct {
	db *sql.DB
}
