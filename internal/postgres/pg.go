package postgres

import (
	"database/sql"
	"fmt"
	"os"

	// DB adapter
	_ "github.com/lib/pq"
)

// Connect ...
func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%s sslmode=%s",
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"))
	fmt.Println("Connecting to DB\n", psqlInfo)
	return sql.Open("postgres", psqlInfo)
}
