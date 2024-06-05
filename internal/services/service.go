package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type services struct {
	userService UserService
	roomService RoomService
	authService AuthService
}

func New(conn *pgxpool.Pool) *services {

	uservice := NewUserService(conn)
	rservice := NewRoomService(conn)
	aservice := NewAuthService(conn)

	return &services{uservice, rservice, aservice}
}

func (s *services) GetUserService() UserService {
	return s.userService
}

func (s *services) GetRoomService() RoomService {
	return s.roomService
}

func (s *services) GetAuthService() AuthService {
	return s.authService
}

type Services interface {
	GetUserService() UserService
	GetRoomService() RoomService
	GetAuthService() AuthService
}
