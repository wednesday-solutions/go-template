-- +migrate Up
create table subjects (
    id SERIAL UNIQUE PRIMARY KEY,
	name text  NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
CREATE INDEX subjects__idx__name ON subjects(name);

-- +migrate Down
drop table subjects;
drop index subjects__idx__name;