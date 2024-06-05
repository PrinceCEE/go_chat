// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email)
VALUES ($1, $2, $3) RETURNING created_at, updated_at, id
`

type CreateUserParams struct {
	FirstName string
	LastName  string
	Email     string
}

type CreateUserRow struct {
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	ID        uuid.UUID
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser, arg.FirstName, arg.LastName, arg.Email)
	var i CreateUserRow
	err := row.Scan(&i.CreatedAt, &i.UpdatedAt, &i.ID)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, first_name, last_name, email, created_at, updated_at FROM users
WHERE id = $1 OR email = $2 LIMIT 1
`

type GetUserParams struct {
	ID    uuid.UUID
	Email string
}

func (q *Queries) GetUser(ctx context.Context, arg GetUserParams) (User, error) {
	row := q.db.QueryRow(ctx, getUser, arg.ID, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, first_name, last_name, email, created_at, updated_at FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET first_name = $1, last_name = $2, email = $3, updated_at = $4
WHERE id = $5
`

type UpdateUserParams struct {
	FirstName string
	LastName  string
	Email     string
	UpdatedAt pgtype.Timestamptz
	ID        uuid.UUID
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}
