package main

import (
	"fmt"
	"go-template/internal/config"
	"go-template/internal/postgres"

	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	migrations := &migrate.FileMigrationSource{
		Dir: "internal/migrations",
	}
	err := config.LoadEnv()
	if err != nil {
		fmt.Println("failed while loading env")
	}
	db, err := postgres.Connect()
	if err != nil {
		fmt.Println("failed while fetching db connection")
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		fmt.Println("failed while executing migration")
	}
	fmt.Printf("Applied %d migrations!\n", n)
}
