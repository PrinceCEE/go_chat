package models

import "time"

type Room struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxMembers  int    `json:"max_members"`
}

type RoomMember struct {
	CreatedAt time.Time `json:"created_at"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
}

type RoomMessage struct {
	UpdatedAt    time.Time `json:"updated_at"`
	RoomID       string    `json:"room_id"`
	RoomMemberID string    `json:"room_member_id"`
	UserID       string    `json:"user_id"`
	Content      string    `json:"content"`
}
