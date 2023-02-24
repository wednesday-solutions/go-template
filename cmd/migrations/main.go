package main

import (
	"database/sql"
	"fmt"
	"go-template/internal/config"
	"go-template/internal/postgres"
	"go-template/pkg/utl/zaplog"
	"log"
	"os"

	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		log.Println("&&&&")
		log.Println(err)
		return
	}
	db, err := postgres.Connect()
	if err != nil {
		fmt.Println("failed while fetching db connection", err)
		zaplog.Logger.Error(err)
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
		zaplog.Logger.Error(err)
		return
	}
	fmt.Printf("Applied %d migrations!\n", n)
}
