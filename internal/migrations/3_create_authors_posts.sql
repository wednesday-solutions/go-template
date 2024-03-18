-- +migrate Up
CREATE TABLE public.authors (
    			id SERIAL UNIQUE PRIMARY KEY,
				first_name TEXT,
				last_name TEXT,
				email TEXT UNIQUE,
				created_at TIMESTAMP WITH TIME ZONE,
				updated_at TIMESTAMP WITH TIME ZONE,
				deleted_at TIMESTAMP WITH TIME ZONE
			);
CREATE INDEX authors_email_idx ON authors(email);
CREATE TABLE public.posts (
                id SERIAL UNIQUE PRIMARY KEY,
                author_id int REFERENCES authors(id) NOT NULL,
                post varchar(255) NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE,
				updated_at TIMESTAMP WITH TIME ZONE,
				deleted_at TIMESTAMP WITH TIME ZONE
			);
CREATE INDEX posts_post_idx ON posts(post);


-- +migrate Down
DROP TABLE posts;
DROP TABLE authors;
