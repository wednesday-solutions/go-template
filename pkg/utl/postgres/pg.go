package postgres

import (
	"database/sql"
	"fmt"
	"github.com/go-pg/pg/v9"
	"log"
	"os"
	// DB adapter
	_ "github.com/lib/pq"
)

// MigrationConnect ...
func MigrationConnect() *pg.DB {
	var psn = os.Getenv("DATABASE_URL")
	u, err := pg.ParseURL(psn)
	if err != nil {
		log.Fatal(err)
	}
	return pg.Connect(u)
}

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
