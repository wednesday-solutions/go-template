-- +migrate Up
CREATE TABLE users (
				id INT AUTO_INCREMENT PRIMARY KEY,
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
				role_id int,
				created_at TIMESTAMP,
				updated_at TIMESTAMP DEFAULT NOW(),
				deleted_at TIMESTAMP,
				CONSTRAINT fk__role_users FOREIGN KEY (role_id) REFERENCES roles(id)
			);
CREATE INDEX username_idx ON users(username);
CREATE INDEX email_idx ON users(email);
CREATE INDEX role_id_idx ON users(role_id);

-- +migrate Down
DROP TABLE users;