package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/go_chat/internal/data/db"
)

type services struct {
	userService *userService
	roomService *roomService
}

func New(conn *pgxpool.Pool) *services {
	store := db.New(conn)

	uservice := &userService{store}
	rservice := &roomService{store}

	return &services{uservice, rservice}
}

func (s *services) GetUserService() UserService {
	return s.userService
}

func (s *services) GetRoomService() RoomService {
	return s.roomService
}

type Services interface {
	GetUserService() UserService
	GetRoomService() RoomService
}
