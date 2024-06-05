package services

import "github.com/princecee/go_chat/internal/data/db"

type roomService struct {
	store *db.Queries
}

func (s *roomService) CreateRoom()     {}
func (s *roomService) GetRoom()        {}
func (s *roomService) GetRooms()       {}
func (s *roomService) DeleteRoom()     {}
func (s *roomService) UpdateRoom()     {}
func (s *roomService) GetRoomMember()  {}
func (s *roomService) LeaveRoom()      {}
func (s *roomService) JoinRoom()       {}
func (s *roomService) GetRoomMembers() {}
func (s *roomService) CreateMessage()  {}
func (s *roomService) GetMessage()     {}
func (s *roomService) GetMessages()    {}
func (s *roomService) DeleteMessage()  {}

type RoomService interface {
	CreateRoom()
	GetRoom()
	GetRooms()
	DeleteRoom()
	UpdateRoom()
	GetRoomMember()
	LeaveRoom()
	JoinRoom()
	GetRoomMembers()
	CreateMessage()
	GetMessage()
	GetMessages()
	DeleteMessage()
}
