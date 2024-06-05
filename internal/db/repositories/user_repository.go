package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *userRepository {
	return &userRepository{conn}
}

func (r *userRepository) CreateUser() {}
func (r *userRepository) GetUser()    {}
func (r *userRepository) GetUsers()   {}
func (r *userRepository) DeleteUser() {}
func (r *userRepository) UpdateUser() {}
