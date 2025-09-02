package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	// go-redis v9
	"github.com/redis/go-redis/v9"
)

var (
	// MySQL
	DB         *sqlx.DB
	EnableAuth bool

	// Redis
	RDB *redis.Client
	Ctx = context.Background()
)

// InitDB initializes MySQL with sane defaults
func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
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
	// connection pool tuning (adjust as needed)
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	EnableAuth, _ = strconv.ParseBool(os.Getenv("ENABLE_AUTH"))
	return nil
}

// InitRedis initializes Redis client
func InitRedis() error {
	addr := getenvDefault("REDIS_ADDR", "127.0.0.1:6379")
	pass := os.Getenv("REDIS_PASSWORD")
	dbn := getenvIntDefault("REDIS_DB", 0)

	RDB = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     pass,
		DB:           dbn,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolSize:     20,
		PoolTimeout:  2 * time.Second,
	})

	ctx, cancel := context.WithTimeout(Ctx, 2*time.Second)
	defer cancel()
	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	log.Println("Connected to Redis:", addr)
	return nil
}

// Helpers / getters
func GetDB() *sqlx.DB               { return DB }
func GetRedisClient() *redis.Client { return RDB }

// Close gracefully closes DB and Redis
func Close() {
	if DB != nil {
		_ = DB.Close()
	}
	if RDB != nil {
		_ = RDB.Close()
	}
}

func getenvDefault(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func getenvIntDefault(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
