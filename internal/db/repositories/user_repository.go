package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	dataSource "github.com/princecee/go_chat/internal/db/data-source"
	"github.com/princecee/go_chat/internal/models"
	"github.com/princecee/go_chat/utils"
)

type userRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *userRepository {
	return &userRepository{conn}
}

func (r *userRepository) CreateUser(user *models.User, tx pgx.Tx) error {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	_user, err := ds.CreateUser(context.Background(), dataSource.CreateUserParams{
		LastName:  user.LastName,
		FirstName: user.FirstName,
		Email:     user.Email,
	})
	if err != nil {
		return err
	}

	user.CreatedAt = _user.CreatedAt
	user.UpdatedAt = _user.UpdatedAt
	user.ID = _user.ID.String()

	return nil
}

type GetUserParams struct {
	ID    string
	Email string
}

func (r *userRepository) GetUser(data GetUserParams, tx pgx.Tx) (*models.User, error) {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	_user, err := ds.GetUser(context.Background(), dataSource.GetUserParams{
		ID:    utils.StringToUUID(data.ID),
		Email: data.Email,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        _user.ID.String(),
		FirstName: _user.FirstName,
		LastName:  _user.LastName,
		Email:     _user.Email,
		CreatedAt: _user.CreatedAt,
		UpdatedAt: _user.UpdatedAt,
	}, nil
}

func (r *userRepository) GetUsers(tx pgx.Tx) ([]*models.User, error) {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	_users, err := ds.GetUsers(context.Background())
	if err != nil {
		return nil, err
	}

	users := []*models.User{}
	for _, _user := range _users {
		user := &models.User{
			ID:        _user.ID.String(),
			FirstName: _user.FirstName,
			LastName:  _user.LastName,
			Email:     _user.Email,
			CreatedAt: _user.CreatedAt,
			UpdatedAt: _user.UpdatedAt,
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) DeleteUser(userId string, tx pgx.Tx) error {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	return ds.DeleteUser(context.Background(), utils.StringToUUID(userId))
}

func (r *userRepository) UpdateUser(user *models.User, tx pgx.Tx) error {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	user.UpdatedAt = time.Now()
	return ds.UpdateUser(context.Background(), dataSource.UpdateUserParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt,
		ID:        utils.StringToUUID(user.ID),
	})
}
