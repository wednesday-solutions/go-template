-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE public.roles (
			id SERIAL UNIQUE PRIMARY KEY,
			access_level int NOT NULL,
			name text  NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE,
			deleted_at TIMESTAMP WITH TIME ZONE
		);


INSERT INTO roles (id,access_level, name, created_at, updated_at, deleted_at) VALUES
(13,1, 'Admin', '2023-06-20 12:00:00+00', '2023-06-20 12:00:00+00', NULL),
(14,2, 'User', '2023-06-20 12:05:00+00', '2023-06-20 12:05:00+00', NULL);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE roles;