-- name: CreateRoom :one
INSERT INTO rooms (name, description, max_members, created_by)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;

-- name: GetRoom :one
SELECT * FROM rooms WHERE id = $1 LIMIT 1;

-- name: GetRooms :many
SELECT * FROM rooms WHERE created_by = COALESCE(sqlc.narg(created_by), created_by);

-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id = $1;

-- name: UpdateRoom :exec
UPDATE rooms SET updated_at = $1, name = $2, description = $3, max_members = $4
WHERE id = $5;

-- name: CreateRoomMember :one
INSERT INTO room_members (room_id, user_id) VALUES($1, $2)
RETURNING id, created_at, updated_at;

-- name: GetRoomMember :one
SELECT * FROM room_members WHERE id = $1 LIMIT 1;

-- name: GetRoomMemberByWhere :one
SELECT * FROM room_members WHERE user_id = $1 AND room_id = $2;

-- name: GetRoomMembers :many
SELECT * FROM room_members WHERE room_id = COALESCE(sqlc.narg(room_id), room_id) AND user_id = COALESCE(sqlc.narg(user_id), user_id);

-- name: DeleteRoomMember :exec
DELETE FROM room_members WHERE id = $1;

-- name: RoomMembersCount :one
SELECT COUNT(*) AS count FROM room_members WHERE room_id = $1;

-- name: CreateRoomMessage :one
INSERT INTO room_messages (room_id, room_member_id, user_id, content)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;

-- name: GetRoomMessage :one
SELECT * FROM room_messages WHERE id = $1 LIMIT 1;

-- name: GetRoomMessages :many
SELECT * FROM room_messages WHERE
  room_id = COALESCE(sqlc.narg(room_id), room_id) AND
  room_member_id = COALESCE(sqlc.narg(room_member_id), room_member_id) AND
  user_id = COALESCE(sqlc.narg(user_id), user_id);

-- name: DeleteRoomMessage :exec
DELETE FROM room_messages WHERE id = $1;