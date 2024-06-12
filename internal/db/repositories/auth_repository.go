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

type authRepository struct {
	conn *pgxpool.Pool
}

func NewAuthRepository(conn *pgxpool.Pool) *authRepository {
	return &authRepository{conn}
}

func (r *authRepository) CreateAuth(auth *models.Auth, tx pgx.Tx) error {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	_auth, err := ds.CreateAuth(context.Background(), dataSource.CreateAuthParams{
		UserID:   utils.StringToUUID(auth.UserID),
		Password: auth.Password,
	})
	if err != nil {
		return err
	}

	auth.CreatedAt = _auth.CreatedAt
	auth.UpdatedAt = _auth.UpdatedAt
	auth.ID = _auth.ID.String()

	return nil
}

func (r *authRepository) GetUserAuth(userId string, tx pgx.Tx) (*models.Auth, error) {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	_auth, err := ds.GetUserAuth(context.Background(), utils.StringToUUID(userId))
	if err != nil {
		return nil, err
	}

	return &models.Auth{
		ID:        _auth.ID.String(),
		CreatedAt: _auth.CreatedAt,
		UpdatedAt: _auth.UpdatedAt,
		UserID:    _auth.UserID.String(),
		Password:  _auth.Password,
	}, nil
}

func (r *authRepository) UpdateUserAuth(auth *models.Auth, tx pgx.Tx) error {
	auth.UpdatedAt = time.Now()

	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	return ds.UpdateUserAuth(context.Background(), dataSource.UpdateUserAuthParams{
		Password:  auth.Password,
		UserID:    utils.StringToUUID(auth.ID),
		UpdatedAt: auth.UpdatedAt,
	})
}
