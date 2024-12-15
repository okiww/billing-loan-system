package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB is the global database connection object
var DB *sql.DB

// Initialize function to set up the MySQL connection
func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
		return err
	}

	// Verify the connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
		return err
	}

	log.Println("Successfully connected to the database")
	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if err := DB.Close(); err != nil {
		log.Fatalf("Error closing the database: %v", err)
		return err
	}
	log.Println("Database connection closed")
	return nil
}

// QueryDB runs a select query and returns the result
func QueryDB(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := DB.Query(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	return rows, nil
}

// ExecuteDB runs an insert/update/delete query
func ExecuteDB(query string, args ...interface{}) (sql.Result, error) {
	result, err := DB.Exec(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	return result, nil
}
