-- name: CreateAuth :one
INSERT INTO auths (user_id, password) VALUES ($1, $2)
RETURNING created_at, updated_at, id;

-- name: GetUserAuth :one
SELECT * FROM auths WHERE user_id = $1 LIMIT 1;

-- name: UpdateUserAuth :exec
UPDATE auths SET password = $1, updated_at = $2 WHERE user_id = $3;