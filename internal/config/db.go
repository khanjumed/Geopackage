package config

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DB         *sqlx.DB
	EnableAuth bool
)

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}

	// Load ENABLE_AUTH from env
	EnableAuth, _ = strconv.ParseBool(os.Getenv("ENABLE_AUTH"))
	return nil
}
