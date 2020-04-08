package main

import (
	"fmt"

	"github.com/go-pg/migrations/v7"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table paytab_users...")

		// creating function for update trigger
		_, err := db.Exec(`CREATE OR REPLACE FUNCTION update_modified_column()
		RETURNS TRIGGER AS $$
		BEGIN
		NEW.modified = now();
		RETURN NEW;
		END;
		$$ language 'plpgsql';`)
		_, err = db.Exec(`
			create table users(
			id SERIAL NOT NULL PRIMARY KEY,
			address TEXT,
			first_name VARCHAR(32) NOT NULL,
			last_name VARCHAR(32) NOT NULL,
			mobile_number NUMERIC,
			email VARCHAR(32) NOT NULL,
			age NUMERIC NOT NULL,
			nationality VARCHAR(32) NOT NULL,
			password VARCHAR(64) NOT NULL,
			token TEXT,
			last_login timestamp(0) default  NULL,
			created_at timestamp(0) default CURRENT_TIMESTAMP NOT NULL,
			updated_at timestamp(0) default CURRENT_TIMESTAMP NOT NULL,
			deleted_at timestamp(0) );
			`)

		// adding update trigger
		_, err = db.Exec(`CREATE TRIGGER update_user_modtime BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table users...")
		_, err := db.Exec(`DROP TABLE users`)
		return err
	})
}
