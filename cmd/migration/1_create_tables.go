package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

type tableQuery struct {
	Name    string
	Query   string
	Columns []string
}

func init() {

	tableQueries := []tableQuery{
		{
			Name: "roles",
			Query: fmt.Sprintf(`CREATE TABLE public.%s (
			id SERIAL UNIQUE PRIMARY KEY,
			access_level int NOT NULL,
			name text  NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
		);`, "roles"),
			Columns: []string{},
		}, {
			Name: "users",
			Query: fmt.Sprintf(`CREATE TABLE public.%s (
				id SERIAL UNIQUE PRIMARY KEY,
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
				role_id int REFERENCES roles(id),
				created_at TIMESTAMP WITH TIME ZONE,
				updated_at TIMESTAMP WITH TIME ZONE,
				deleted_at TIMESTAMP WITH TIME ZONE
			);`, "users"),
			Columns: []string{"username", "email", "role_id"},
		},
	}

	migrations.MustRegister(func(db migrations.DB) error {
		for i := 0; i < len(tableQueries); i++ {
			err := CreateTriggerForUpdatedAt(db)
			if err != nil {
				return err
			}
			err = CreateTableAndIndexes(db, tableQueries[i].Query, tableQueries[i].Name, tableQueries[i].Columns)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		for i := len(tableQueries) - 1; i >= 0; i-- {
			err := DropTable(db, tableQueries[i].Name)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
