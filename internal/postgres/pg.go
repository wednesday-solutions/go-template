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
		os.Getenv("PSQL_DBNAME"),
		os.Getenv("PSQL_HOST"),
		os.Getenv("PSQL_USER"),
		os.Getenv("PSQL_PASS"),
		os.Getenv("PSQL_PORT"),
		os.Getenv("PSQL_SSLMODE"))
	fmt.Println("Connecting to DB\n", psqlInfo)
	return sql.Open("postgres", psqlInfo)
}
