package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
    _ "github.com/go-sql-driver/mysql"
)

var db *sqlx.DB

func InitDB() error {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	var err error
	db, err = sqlx.Connect("mysql", dsn)
	return err
}
