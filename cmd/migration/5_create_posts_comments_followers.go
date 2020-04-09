package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {

	var tableName = "posts"
	var createTableQuery = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(200) NOT NULL,
			body TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
			);`, tableName)
	migrations.MustRegister(func(db migrations.DB) error {
		err := createTable(db, createTableQuery, tableName)
		if err != nil {
			return err
		}
		tableName = "comments"
		createTableQuery = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			title VARCHAR(200) NOT NULL,
			body TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
			);`, tableName)
		err = createTable(db, createTableQuery, tableName)

		tableName = "followers"
		createTableQuery = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s (
			follower_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			followee_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE,
			PRIMARY KEY(follower_id, followee_id)
			);`, tableName)
		err = createTable(db, createTableQuery, tableName)
		return err
	}, func(db migrations.DB) error {
		return DropTable(db, tableName)
	})
}

func createTable(db migrations.DB, createTableQuery string, tableName string) error {
	err := CreateTableAndAddTrigger(db, createTableQuery, tableName)
	if err != nil {
		return err
	}
	return err
}
