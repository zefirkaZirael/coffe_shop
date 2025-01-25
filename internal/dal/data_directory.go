package dal

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Connects to the Database
func CheckDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
