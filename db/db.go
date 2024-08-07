package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func DBConnect() error {
	username := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")
	database := os.Getenv("PG_DATABASE")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")

	// Construct connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, host, port, database)

	// Connect to the PostgreSQL database
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return err
	}
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(time.Hour)
	err = DB.Ping()
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return err
	}
	return nil

}
func Conntest() {
	username := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")
	database := os.Getenv("PG_DATABASE")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")

	// Construct connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, host, port, database)

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return
	}
	defer db.Close()

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: ", err)
		return
	}
	fmt.Println("Hello!")

}
