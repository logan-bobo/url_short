// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package database

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING id, email, password, created_at, updated_at
`

type CreateUserParams struct {
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.Password,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const selectUser = `-- name: SelectUser :one
SELECT id, email, password, created_at, updated_at
FROM users
WHERE email = $1
`

func (q *Queries) SelectUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, selectUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const selectUserByID = `-- name: SelectUserByID :one
SELECT id, email, password, created_at, updated_at
FROM users
WHERE id = $1
`

func (q *Queries) SelectUserByID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, selectUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET email = $1, password = $2, updated_at = $3
WHERE id = $4
`

type UpdateUserParams struct {
	Email     string
	Password  string
	UpdatedAt time.Time
	ID        int32
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.Email,
		arg.Password,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}
