package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type services struct {
	userService UserService
	roomService RoomService
	authService AuthService
	conn        *pgxpool.Pool
}

var _services *services

func New(conn *pgxpool.Pool) *services {
	if _services != nil {
		return _services
	}

	uservice := NewUserService(conn)
	rservice := NewRoomService(conn)
	aservice := NewAuthService(conn)

	return &services{uservice, rservice, aservice, conn}
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

func (s *services) GetDB() *pgxpool.Pool {
	return s.conn
}

type Services interface {
	GetUserService() UserService
	GetRoomService() RoomService
	GetAuthService() AuthService
	GetDB() *pgxpool.Pool
}
