package models

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/antonve/logger-api/config"
)

const maxOpenConns = 8
const maxIdleConns = 8

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

	sqlxDB.SetMaxOpenConns(maxOpenConns)
	sqlxDB.SetMaxIdleConns(maxIdleConns)

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

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)

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

	sqlxConnection.SetMaxOpenConns(maxOpenConns)
	sqlxConnection.SetMaxIdleConns(maxIdleConns)

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

	sqlConnection.SetMaxOpenConns(maxOpenConns)
	sqlConnection.SetMaxIdleConns(maxIdleConns)

	return sqlConnection
}
