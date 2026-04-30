package utils

import (
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func Db_Init() (*sqlx.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", db_user, db_password, db_name)

	var db *sqlx.DB

	for i := 0; i < 10; i++ {
		db, err = sqlx.Connect("pgx", connectionString)
		if err == nil {
			return db, nil
		}

		fmt.Println("waiting for database...")
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after retries: %w", err)
}
