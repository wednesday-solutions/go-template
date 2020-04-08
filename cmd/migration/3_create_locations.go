package main

import (
	"fmt"
	"github.com/go-pg/migrations/v7"
)

func init() {
	var tableName = "locations"
	var createTableQuery = fmt.Sprintf(`CREATE TABLE public.%s (
			id serial unique,
			created_at timestamp with time zone,
			updated_at timestamp with time zone,
			deleted_at timestamp with time zone,
			name text,
			active boolean,
			address text,
			company_id serial REFERENCES companies(id)
		);`, tableName)
	migrations.MustRegister(func(db migrations.DB) error {
		err := CreateTriggerForUpdatedAt(db)
		if err != nil {
			return err
		}
		err = CreateTableAndAddTrigger(db, createTableQuery, tableName)
		if err != nil {
			return err
		}
		return err
	}, func(db migrations.DB) error {
		return DropTable(db, tableName)
	})
}
