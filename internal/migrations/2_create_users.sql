-- +migrate Up
CREATE TABLE users (
				id SERIAL UNIQUE PRIMARY KEY,
				first_name TEXT,
				last_name TEXT,
				username VARCHAR(100) UNIQUE,
				password TEXT,
				email VARCHAR(100) UNIQUE,
				mobile TEXT,
				address TEXT,
				active BOOLEAN,
				last_login TIMESTAMP ,
				last_password_change TIMESTAMP,
				token TEXT,
				role_id int REFERENCES roles(id),
				created_at TIMESTAMP,
				updated_at TIMESTAMP,
				deleted_at TIMESTAMP
			);
CREATE INDEX username_idx ON users(username);
CREATE INDEX email_idx ON users(email);
CREATE INDEX role_id_idx ON users(role_id);

-- +migrate Down
DROP TABLE users;