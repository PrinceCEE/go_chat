package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/go_chat/internal/db/repositories"
)

type roomService struct {
	conn           *pgxpool.Pool
	RoomRepository RoomRepository
}

func NewRoomService(conn *pgxpool.Pool) RoomService {
	return &roomService{
		conn:           conn,
		RoomRepository: repositories.NewRoomRepository(conn),
	}
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

type RoomRepository interface {
	CreateRoom()
	GetRoom()
	GetRooms()
	DeleteRoom()
	UpdateRoom()
	CreateRoomMember()
	GetRoomMember()
	GetRoomMembers()
	DeleteRoomMember()
	UpdateRoomMember()
	CreateRoomMessage()
	GetRoomMessage()
	GetRoomMessages()
	DeleteRoomMessage()
}

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
