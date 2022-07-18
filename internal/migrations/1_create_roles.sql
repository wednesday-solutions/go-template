-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE roles (
			id INT AUTO_INCREMENT PRIMARY KEY,
			access_level int NOT NULL,
			name text  NOT NULL,
			created_at TIMESTAMP,
			updated_at TIMESTAMP,
			deleted_at TIMESTAMP
		);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE roles;