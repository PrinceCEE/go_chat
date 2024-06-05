package models

type Room struct {
	ModelMixin
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxMembers  int    `json:"max_members"`
}

type RoomMember struct {
	ModelMixin
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

type RoomMessage struct {
	ModelMixin
	RoomID       string `json:"room_id"`
	RoomMemberID string `json:"room_member_id"`
	UserID       string `json:"user_id"`
	Content      string `json:"content"`
}
