-- +migrate Up
CREATE TABLE public.users (
				id SERIAL UNIQUE PRIMARY KEY,
				first_name TEXT,
				last_name TEXT,
				username TEXT UNIQUE,
				password TEXT,
				email TEXT UNIQUE,
				mobile TEXT,
				address TEXT,
				active BOOLEAN,
				last_login TIMESTAMP WITH TIME ZONE,
				last_password_change TIMESTAMP WITH TIME ZONE,
				token TEXT,
				role_id int REFERENCES roles(id),
				created_at TIMESTAMP WITH TIME ZONE,
				updated_at TIMESTAMP WITH TIME ZONE,
				deleted_at TIMESTAMP WITH TIME ZONE
			);
CREATE INDEX users_username_idx ON users(username);
CREATE INDEX users_email_idx ON users(email);
CREATE INDEX users_role_id_idx ON users(role_id);

-- +migrate Down
DROP TABLE users;