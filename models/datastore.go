package models

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/antonve/logger-api/config"
)

var sqlxDB *sqlx.DB
var sqlDB *sql.DB
var sqlxConnection *sqlx.DB
var sqlConnection *sql.DB

// GetDatabase Returns database connection from pool
func GetDatabase() *sqlx.DB {
	if sqlxDB != nil {
		return sqlxDB
	}

	sqlxDB, err := sqlx.Open("postgres", config.GetConfig().GetCompleteConnectionString())
	if err != nil {
		log.Fatalln("Couldn't connect to data store")

		return nil
	}

	return sqlxDB
}

// GetSQLDatabase Returns database connection from pool with the default sql package
func GetSQLDatabase() *sql.DB {
	if sqlDB != nil {
		return sqlDB
	}

	sqlDB, err := sql.Open("postgres", config.GetConfig().GetCompleteConnectionString())
	if err != nil {
		log.Fatalln("Couldn't connect to data store")

		return nil
	}

	return sqlDB
}

// GetConnection Returns database connection without a database
func GetConnection() *sqlx.DB {
	if sqlxConnection != nil {
		return sqlxConnection
	}

	sqlxConnection, err := sqlx.Open("postgres", config.GetConfig().ConnectionString)
	if err != nil {
		log.Fatalln("Couldn't connect to data store")

		return nil
	}

	return sqlxConnection
}

// GetSQLConnection Returns database connection without a database with the default sql package
func GetSQLConnection() *sql.DB {
	if sqlConnection != nil {
		return sqlConnection
	}

	sqlConnection, err := sql.Open("postgres", config.GetConfig().ConnectionString)
	if err != nil {
		log.Fatalln("Couldn't connect to data store")

		return nil
	}

	return sqlConnection
}
