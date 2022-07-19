package ms

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// ...
func Connect() (*sql.DB, error) {
	return sql.Open("mysql", "root:password@/erp_local")
}
