-- +goose Up
CREATE TABLE urls (
		id integer PRIMARY KEY,
		short_url VARCHAR(100) UNIQUE NOT NULL,
		long_url VARCHAR(300) UNIQUE NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE urls;
