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

		fmt.Print("Creating users table")
		_, err = db.Exec(`
			CREATE TABLE public.users (
				id bigint NOT NULL,
				created_at timestamp with time zone,
				updated_at timestamp with time zone,
				deleted_at timestamp with time zone,
				first_name text,
				last_name text,
				username text,
				password text,
				email text,
				mobile text,
				phone text,
				address text,
				active boolean,
				last_login timestamp with time zone,
				last_password_change timestamp with time zone,
				token text,
				role_id bigint,
				company_id bigint,
				location_id bigint
			);`)

		// adding update trigger
		_, err = db.Exec(`CREATE TRIGGER update_user_modtime BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();`)
		fmt.Print("Done creating users table\n Creating locations table")

		db.Exec(`CREATE TABLE public.locations (
			id bigint NOT NULL,
			created_at timestamp with time zone,
			updated_at timestamp with time zone,
			deleted_at timestamp with time zone,
			name text,
			active boolean,
			address text,
			company_id bigint
		);`)

		_, err = db.Exec(`CREATE TRIGGER update_user_modtime BEFORE UPDATE ON locations FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();`)

		fmt.Print("Done creating users table\n Creating roles table")

		_, err = db.Exec(`CREATE TABLE public.roles (
			id bigint NOT NULL,
			access_level bigint,
			name text
		);`)

		_, err = db.Exec(`CREATE TRIGGER update_user_modtime BEFORE UPDATE ON roles FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();`)

		fmt.Print("Done creating roles table\n Creating companies table")

		_, err = db.Exec(`CREATE TABLE public.companies (
			id bigint NOT NULL,
			created_at timestamp with time zone,
			updated_at timestamp with time zone,
			deleted_at timestamp with time zone,
			name text,
			active boolean
		);`)

		_, err = db.Exec(`CREATE TRIGGER update_user_modtime BEFORE UPDATE ON companies FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table users...")
		_, err := db.Exec(`DROP TABLE users`)
		return err
	})
}
