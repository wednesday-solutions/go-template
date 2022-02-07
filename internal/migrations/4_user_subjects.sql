-- +migrate Up
create table user_subjects (
    id SERIAL UNIQUE PRIMARY KEY,
	subject_id int REFERENCES subjects(id),
    user_id int REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
create index user_subjects__idx__subject_id on user_subjects(subject_id);
create index user_subjects__idx__user_id on user_subjects(user_id);

-- +migrate Down
drop table user_subjects;
drop index user_subjects__idx__subject_id;
drop index user_subjects__idx__user_id;