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

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE roles;