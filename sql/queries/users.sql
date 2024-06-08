-- name: CreateUser :one
INSERT INTO users (email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: SelectUser :one
SELECT *
FROM users
WHERE email = $1;

-- name: SelectUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET email = $1, password = $2, updated_at = $3
WHERE id = $4;
