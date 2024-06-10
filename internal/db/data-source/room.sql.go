// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: room.sql

package dataSource

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createRoom = `-- name: CreateRoom :one
INSERT INTO rooms (name, description, max_members, created_by)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at
`

type CreateRoomParams struct {
	Name        string
	Description pgtype.Text
	MaxMembers  int32
	CreatedBy   uuid.UUID
}

type CreateRoomRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) (CreateRoomRow, error) {
	row := q.db.QueryRow(ctx, createRoom,
		arg.Name,
		arg.Description,
		arg.MaxMembers,
		arg.CreatedBy,
	)
	var i CreateRoomRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const createRoomMember = `-- name: CreateRoomMember :one
INSERT INTO room_members (room_id, user_id) VALUES($1, $2)
RETURNING id, created_at, updated_at
`

type CreateRoomMemberParams struct {
	RoomID uuid.UUID
	UserID uuid.UUID
}

type CreateRoomMemberRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateRoomMember(ctx context.Context, arg CreateRoomMemberParams) (CreateRoomMemberRow, error) {
	row := q.db.QueryRow(ctx, createRoomMember, arg.RoomID, arg.UserID)
	var i CreateRoomMemberRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const createRoomMessage = `-- name: CreateRoomMessage :one
INSERT INTO room_messages (room_id, room_member_id, user_id, content)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at
`

type CreateRoomMessageParams struct {
	RoomID       uuid.UUID
	RoomMemberID uuid.UUID
	UserID       uuid.UUID
	Content      string
}

type CreateRoomMessageRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateRoomMessage(ctx context.Context, arg CreateRoomMessageParams) (CreateRoomMessageRow, error) {
	row := q.db.QueryRow(ctx, createRoomMessage,
		arg.RoomID,
		arg.RoomMemberID,
		arg.UserID,
		arg.Content,
	)
	var i CreateRoomMessageRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deleteRoom = `-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id = $1
`

func (q *Queries) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteRoom, id)
	return err
}

const deleteRoomMember = `-- name: DeleteRoomMember :exec
DELETE FROM room_members WHERE id = $1
`

func (q *Queries) DeleteRoomMember(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteRoomMember, id)
	return err
}

const deleteRoomMessage = `-- name: DeleteRoomMessage :exec
DELETE FROM room_messages WHERE id = $1
`

func (q *Queries) DeleteRoomMessage(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteRoomMessage, id)
	return err
}

const getRoom = `-- name: GetRoom :one
SELECT id, name, description, max_members, created_by, created_at, updated_at FROM rooms WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRoom(ctx context.Context, id uuid.UUID) (Room, error) {
	row := q.db.QueryRow(ctx, getRoom, id)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.MaxMembers,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRoomMember = `-- name: GetRoomMember :one
SELECT id, room_id, user_id, created_at, updated_at FROM room_members WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRoomMember(ctx context.Context, id uuid.UUID) (RoomMember, error) {
	row := q.db.QueryRow(ctx, getRoomMember, id)
	var i RoomMember
	err := row.Scan(
		&i.ID,
		&i.RoomID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRoomMembers = `-- name: GetRoomMembers :many
SELECT id, room_id, user_id, created_at, updated_at FROM room_members WHERE room_id = COALESCE($1, room_id) AND user_id = COALESCE($2, user_id)
`

type GetRoomMembersParams struct {
	RoomID pgtype.UUID
	UserID pgtype.UUID
}

func (q *Queries) GetRoomMembers(ctx context.Context, arg GetRoomMembersParams) ([]RoomMember, error) {
	rows, err := q.db.Query(ctx, getRoomMembers, arg.RoomID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RoomMember
	for rows.Next() {
		var i RoomMember
		if err := rows.Scan(
			&i.ID,
			&i.RoomID,
			&i.UserID,
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

const getRoomMessage = `-- name: GetRoomMessage :one
SELECT id, room_id, room_member_id, user_id, content, created_at, updated_at FROM room_messages WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRoomMessage(ctx context.Context, id uuid.UUID) (RoomMessage, error) {
	row := q.db.QueryRow(ctx, getRoomMessage, id)
	var i RoomMessage
	err := row.Scan(
		&i.ID,
		&i.RoomID,
		&i.RoomMemberID,
		&i.UserID,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRoomMessages = `-- name: GetRoomMessages :many
SELECT id, room_id, room_member_id, user_id, content, created_at, updated_at FROM room_messages WHERE
  room_id = COALESCE($1, room_id) AND
  room_member_id = COALESCE($2, room_member_id) AND
  user_id = COALESCE($3, user_id)
`

type GetRoomMessagesParams struct {
	RoomID       pgtype.UUID
	RoomMemberID pgtype.UUID
	UserID       pgtype.UUID
}

func (q *Queries) GetRoomMessages(ctx context.Context, arg GetRoomMessagesParams) ([]RoomMessage, error) {
	rows, err := q.db.Query(ctx, getRoomMessages, arg.RoomID, arg.RoomMemberID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RoomMessage
	for rows.Next() {
		var i RoomMessage
		if err := rows.Scan(
			&i.ID,
			&i.RoomID,
			&i.RoomMemberID,
			&i.UserID,
			&i.Content,
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

const getRooms = `-- name: GetRooms :many
SELECT id, name, description, max_members, created_by, created_at, updated_at FROM rooms WHERE created_by = COALESCE($1, created_by)
`

func (q *Queries) GetRooms(ctx context.Context, createdBy pgtype.UUID) ([]Room, error) {
	rows, err := q.db.Query(ctx, getRooms, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Room
	for rows.Next() {
		var i Room
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.MaxMembers,
			&i.CreatedBy,
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

const updateRoom = `-- name: UpdateRoom :exec
UPDATE rooms SET updated_at = $1, name = $2, description = $3, max_members = $4
WHERE id = $5
`

type UpdateRoomParams struct {
	UpdatedAt   time.Time
	Name        string
	Description pgtype.Text
	MaxMembers  int32
	ID          uuid.UUID
}

func (q *Queries) UpdateRoom(ctx context.Context, arg UpdateRoomParams) error {
	_, err := q.db.Exec(ctx, updateRoom,
		arg.UpdatedAt,
		arg.Name,
		arg.Description,
		arg.MaxMembers,
		arg.ID,
	)
	return err
}