package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	var tableName = "users"
	var createTableQuery = fmt.Sprintf(`CREATE TABLE public.%s (
				id SERIAL UNIQUE PRIMARY KEY,
				created_at TIMESTAMP WITH TIME ZONE,
				updated_at TIMESTAMP WITH TIME ZONE,
				deleted_at TIMESTAMP WITH TIME ZONE,
				first_name TEXT,
				last_name TEXT,
				username TEXT UNIQUE,
				password TEXT,
				email TEXT UNIQUE,
				mobile TEXT,
				phone TEXT,
				address TEXT,
				active BOOLEAN,
				last_login TIMESTAMP WITH TIME ZONE,
				last_password_change TIMESTAMP WITH TIME ZONE,
				token TEXT,
				role_id BIGINT REFERENCES roles(id),
				company_id BIGINT REFERENCES companies(id),
				location_id BIGINT REFERENCES locations(id)
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
