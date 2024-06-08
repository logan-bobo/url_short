-- +goose Up 
CREATE TABLE users (
	id serial PRIMARY KEY,
	email VARCHAR(250) UNIQUE NOT NULL,
	password VARCHAR(250) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	refresh_token varchar(250),
	refresh_token_revoke_date TIMESTAMP
);

-- +goose Down
DROP TABLE users;
