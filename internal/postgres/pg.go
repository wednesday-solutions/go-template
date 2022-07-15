package postgres

import (
	"database/sql"
	"fmt"
	"os"

	// DB adapter
	_ "github.com/go-sql-driver/mysql"
)

// Connect ...
func Connect() (*sql.DB, error) {
	mysql := fmt.Sprintf("user=%spassword@/dbname=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"))
	fmt.Println("Connecting to DB\n", mysql)
	return sql.Open("mysql", mysql)
}
