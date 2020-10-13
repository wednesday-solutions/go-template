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
			Name: "companies",
			Query: fmt.Sprintf(`CREATE TABLE public.%s (
			id SERIAL UNIQUE PRIMARY KEY,
			name text,
			active boolean,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
		);`, "companies"),
			Columns: []string{},
		}, {
			Name: "locations",
			Query: fmt.Sprintf(`CREATE TABLE public.%s (
			id SERIAL UNIQUE PRIMARY KEY,
			name TEXT,
			active BOOLEAN,
			address TEXT,
			company_id SERIAL REFERENCES companies(id),
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
		);`, "locations"),
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
				company_id int REFERENCES companies(id),
				location_id int REFERENCES locations(id),
				created_at TIMESTAMP WITH TIME ZONE,
				updated_at TIMESTAMP WITH TIME ZONE,
				deleted_at TIMESTAMP WITH TIME ZONE
			);`, "users"),
			Columns: []string{"username"},
		}, {
			Name: "posts",
			Query: fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(200) NOT NULL,
			body TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
			);`, "posts"),
			Columns: []string{},
		}, {
			Name: "comments",
			Query: fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s (	
			id SERIAL PRIMARY KEY,	
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,	
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,	
			title VARCHAR(200) NOT NULL,	
			body TEXT NOT NULL,	
			created_at TIMESTAMP WITH TIME ZONE,	
			updated_at TIMESTAMP WITH TIME ZONE,	
			deleted_at TIMESTAMP WITH TIME ZONE	
			);`, "comments"),
			Columns: []string{},
		}, {
			Name: "followers",
			Query: fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s (
			id SERIAL NOT NULL PRIMARY KEY,
			follower_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			followee_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
			);`, "followers"),
			Columns: []string{},
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
