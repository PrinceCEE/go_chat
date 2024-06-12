package services

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/go_chat/internal/db/repositories"
	"github.com/princecee/go_chat/internal/models"
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

func (s *userService) CreateUser(user *models.User, tx pgx.Tx) error {
	return s.UserRepository.CreateUser(user, tx)
}

func (s *userService) GetUser(data repositories.GetUserParams, tx pgx.Tx) (*models.User, error) {
	return s.UserRepository.GetUser(data, tx)
}

func (s *userService) GetUsers(tx pgx.Tx) ([]*models.User, error) {
	return s.UserRepository.GetUsers(tx)
}

func (s *userService) DeleteUser(userId string, tx pgx.Tx) error {
	return s.UserRepository.DeleteUser(userId, tx)
}

func (s *userService) UpdateUser(user *models.User, tx pgx.Tx) error {
	return s.UserRepository.UpdateUser(user, tx)
}

type UserService interface {
	CreateUser(user *models.User, tx pgx.Tx) error
	GetUser(data repositories.GetUserParams, tx pgx.Tx) (*models.User, error)
	GetUsers(tx pgx.Tx) ([]*models.User, error)
	DeleteUser(userId string, tx pgx.Tx) error
	UpdateUser(user *models.User, tx pgx.Tx) error
}

type UserRepository interface {
	CreateUser(user *models.User, tx pgx.Tx) error
	GetUser(data repositories.GetUserParams, tx pgx.Tx) (*models.User, error)
	GetUsers(tx pgx.Tx) ([]*models.User, error)
	DeleteUser(userId string, tx pgx.Tx) error
	UpdateUser(user *models.User, tx pgx.Tx) error
}
