-- +goose Up
ALTER TABLE urls
ADD COLUMN user_id int NOT NULL,
ADD CONSTRAINT fk_user
	FOREIGN KEY (user_ID)
		REFERENCES users(id)
		    ON DELETE CASCADE;

-- +goose Down
ALTER TABLE urls
DROP COLUMN user_id;

