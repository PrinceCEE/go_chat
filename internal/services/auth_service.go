package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/go_chat/internal/db/repositories"
	"github.com/princecee/go_chat/internal/models"
)

type authService struct {
	conn           *pgxpool.Pool
	AuthRepository AuthRepository
}

func NewAuthService(conn *pgxpool.Pool) AuthService {
	return &authService{
		conn:           conn,
		AuthRepository: repositories.NewAuthRepository(conn),
	}
}

func (s *authService) CreateAuth(auth *models.Auth, tx *pgxpool.Tx) error {
	return s.AuthRepository.CreateAuth(auth, tx)
}

func (s *authService) GetUserAuth(userId string, tx *pgxpool.Tx) (*models.Auth, error) {
	return s.AuthRepository.GetUserAuth(userId, tx)
}

func (s *authService) UpdateUserAuth(auth *models.Auth, tx *pgxpool.Tx) error {
	return s.AuthRepository.UpdateUserAuth(auth, tx)
}

type AuthRepository interface {
	CreateAuth(auth *models.Auth, tx *pgxpool.Tx) error
	GetUserAuth(userId string, tx *pgxpool.Tx) (*models.Auth, error)
	UpdateUserAuth(auth *models.Auth, tx *pgxpool.Tx) error
}

type AuthService interface {
	CreateAuth(auth *models.Auth, tx *pgxpool.Tx) error
	GetUserAuth(userId string, tx *pgxpool.Tx) (*models.Auth, error)
	UpdateUserAuth(auth *models.Auth, tx *pgxpool.Tx) error
}
