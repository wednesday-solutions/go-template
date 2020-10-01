package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {

	tableNames := []string{"roles", "companies", "locations", "users", "posts", "comments", "followers"}

	var createTableQueries = []string{
		fmt.Sprintf(`CREATE TABLE public.%s (
			id SERIAL UNIQUE PRIMARY KEY,
			access_level int NOT NULL,
			name text  NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
		);`, tableNames[0]),
		fmt.Sprintf(`CREATE TABLE public.%s (
			id SERIAL UNIQUE PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE,
			name text,
			active boolean
		);`, tableNames[1]),
		fmt.Sprintf(`CREATE TABLE public.%s (
			id SERIAL UNIQUE PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE,
			name TEXT,
			active BOOLEAN,
			address TEXT,
			company_id SERIAL REFERENCES companies(id)
		);`, tableNames[2]),
		fmt.Sprintf(`CREATE TABLE public.%s (
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
				role_id int REFERENCES roles(id),
				company_id int REFERENCES companies(id),
				location_id int REFERENCES locations(id)
			);`, tableNames[3]),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(200) NOT NULL,
			body TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
			);`, tableNames[4]),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s (	
			id SERIAL PRIMARY KEY,	
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,	
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,	
			title VARCHAR(200) NOT NULL,	
			body TEXT NOT NULL,	
			created_at TIMESTAMP WITH TIME ZONE,	
			updated_at TIMESTAMP WITH TIME ZONE,	
			deleted_at TIMESTAMP WITH TIME ZONE	
			);`, tableNames[5]),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s (
			id SERIAL NOT NULL PRIMARY KEY,
			follower_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			followee_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
			);`, tableNames[6]),
	}

	migrations.MustRegister(func(db migrations.DB) error {
		for i := 0; i < len(tableNames); i++ {
			err := CreateTriggerForUpdatedAt(db)
			if err != nil {
				return err
			}
			err = createTable(db, createTableQueries[i], tableNames[i])
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		for i := len(tableNames) - 1; i >= 0; i-- {
			err := DropTable(db, tableNames[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func createTable(db migrations.DB, createTableQuery string, tableName string) error {
	err := CreateTableAndAddTrigger(db, createTableQuery, tableName)
	if err != nil {
		return err
	}
	return err
}
