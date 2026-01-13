package db

import (
	"chat-application/internal/constants"
	"chat-application/util"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabase() (*sql.DB, error) {
	env := util.GetEnv("ENVIRONMENT", "development")

	var db *sql.DB
	var err error

	if env != "production" {
		dbHost := util.GetEnv("DB_HOST", "localhost")
		dbPort := util.GetEnv("DB_PORT", "5432")
		dbUser := util.GetEnv("DB_USER", "postgres")
		dbPassword := util.GetEnv("DB_PASSWORD", "password")
		dbName := util.GetEnv("DB_NAME", "chat_app")

		localDSN := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName,
		)

		log.Println("Connecting to local database...")
		log.Printf("Environment: %s", env)
		log.Printf("Connecting to: %s:%s/%s", dbHost, dbPort, dbName)

		db, err = sql.Open("pgx", localDSN)
		if err != nil {
			log.Fatalf("Failed to connect to local database: %v", err)
		}
	} else {
		connStr := util.GetEnv("DATABASE_URL", "")
		if connStr == "" {
			log.Fatal("DATABASE_URL is not set for production environment")
		}

		log.Printf("Environment: %s", env)
		log.Println("Connecting to production database...")

		db, err = sql.Open("pgx", connStr)
		if err != nil {
			log.Fatalf("Failed to connect to production database: %v", err)
		}
	}

	// Configure connection pool
	db.SetMaxOpenConns(constants.DBMaxOpenConns)
	db.SetMaxIdleConns(constants.DBMaxIdleConns)
	db.SetConnMaxLifetime(constants.DBConnMaxLifetime)

	return db, nil
}
