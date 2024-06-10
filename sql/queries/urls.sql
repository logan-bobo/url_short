-- name: CreateURL :one
INSERT INTO urls (short_url, long_url, created_at, updated_at, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: SelectURL :one
SELECT * 
FROM urls
WHERE short_url = $1;
