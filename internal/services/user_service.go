package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/go_chat/internal/db/repositories"
)

type userService struct {
	conn           *pgxpool.Pool
	UserRepository UserRepository
}

func NewUserService(conn *pgxpool.Pool) UserService {
	return &userService{
		conn:           conn,
		UserRepository: repositories.NewUserRepository(conn),
	}
}

func (s *userService) CreateUser() {}

func (s *userService) GetUser() {}

func (s *userService) GetUsers() {}

func (s *userService) DeleteUser() {}

func (s *userService) UpdateUser() {}

type UserService interface {
	CreateUser()
	GetUser()
	GetUsers()
	DeleteUser()
	UpdateUser()
}

type UserRepository interface {
	CreateUser()
	GetUser()
	GetUsers()
	DeleteUser()
	UpdateUser()
}
