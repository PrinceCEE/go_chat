-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email)
VALUES ($1, $2, $3) RETURNING created_at, updated_at, id;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 OR email = $2 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET first_name = $1, last_name = $2, email = $3, updated_at = $4
WHERE id = $5;