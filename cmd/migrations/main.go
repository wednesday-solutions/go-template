package main

import (
	"database/sql"
	"fmt"
	"go-template/internal/config"
	"go-template/internal/postgres"
	"os"

	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		fmt.Println("failed while loading env")
	}
	db, err := postgres.Connect()
	if err != nil {
		fmt.Println("failed while fetching db connection")
	}

	for _, arg := range os.Args {
		if arg == "down" {
			runMigration(db, migrate.Down)
		}
	}
	runMigration(db, migrate.Up)
}

func runMigration(db *sql.DB, direction migrate.MigrationDirection) {
	migrations := &migrate.FileMigrationSource{
		Dir: "internal/migrations",
	}
	n, err := migrate.Exec(db, "postgres", migrations, direction)
	if err != nil {
		fmt.Println("failed while executing migration")
	}
	fmt.Printf("Applied %d migrations!\n", n)
}
