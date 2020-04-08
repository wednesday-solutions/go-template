package main

import (
	"fmt"
	"github.com/go-pg/migrations/v7"
)

func init() {
	var tableName = "users"
	var createTableQuery = fmt.Sprintf(`CREATE TABLE public.%s (
				id serial unique,
				created_at timestamp with time zone,
				updated_at timestamp with time zone,
				deleted_at timestamp with time zone,
				first_name text,
				last_name text,
				username text unique,
				password text,
				email text unique,
				mobile text,
				phone text,
				address text,
				active boolean,
				last_login timestamp with time zone,
				last_password_change timestamp with time zone,
				token text,
				role_id bigint REFERENCES roles(id),
				company_id bigint REFERENCES companies(id),
				location_id bigint REFERENCES locations(id)
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
