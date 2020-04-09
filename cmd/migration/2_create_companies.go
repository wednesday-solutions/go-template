package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	var tableName = "companies"
	var createTableQuery = fmt.Sprintf(`CREATE TABLE public.%s (
			id SERIAL UNIQUE PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE,
			name text,
			active boolean
		);`, tableName)
	migrations.MustRegister(func(db migrations.DB) error {
		err := CreateTableAndAddTrigger(db, createTableQuery, tableName)
		if err != nil {
			return err
		}
		return err
	}, func(db migrations.DB) error {
		return DropTable(db, tableName)
	})
}
