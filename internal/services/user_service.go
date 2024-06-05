package services

import "github.com/princecee/go_chat/internal/data/db"

type userService struct {
	store *db.Queries
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
