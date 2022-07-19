package mysql

import (
	"database/sql"
	"fmt"
	"os"

	// DB adapter
	_ "github.com/go-sql-driver/mysql"
)

// Connect ...
func Connect() (*sql.DB, error) {
	mysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_ROOT_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DBNAME"))
	fmt.Println("Connecting to DB\n", mysql)
	return sql.Open("mysql", mysql)
}
