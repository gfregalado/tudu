package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

func Init(dotenvPath string) (*sql.DB, error) {
	err := godotenv.Load(dotenvPath)
	if err != nil {
		return nil, err
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dbURI := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPass, dbPort, dbName)

	db, err := sql.Open("postgres", dbURI)

	/*defer db.Close()*/

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
