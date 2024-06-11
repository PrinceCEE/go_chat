package services

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/princecee/go_chat/internal/db/repositories"
	"github.com/princecee/go_chat/internal/models"
	"github.com/princecee/go_chat/utils"
	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	suite.Suite
	conn        *pgxpool.Pool
	authService *authService
	userService *userService
}

func (s *AuthServiceTestSuite) SetupSuite() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}

	conn, err := pgxpool.New(context.Background(), os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}

	authRepository := repositories.NewAuthRepository(conn)
	userRepository := repositories.NewUserRepository(conn)
	s.conn = conn
	s.authService = &authService{conn: conn, AuthRepository: authRepository}
	s.userService = &userService{conn: conn, UserRepository: userRepository}
}

func (s *AuthServiceTestSuite) TearDownSuite() {
	defer s.conn.Close()

	teardownQuery := `
		DELETE FROM users;
		DELETE FROM rooms;
	`

	_, err := s.conn.Exec(context.Background(), teardownQuery)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			panic(err)
		}
	}
}

func (s *AuthServiceTestSuite) TestAuthService() {
	user := models.User{
		FirstName: "Chimezie",
		LastName:  "Edeh",
		Email:     "princecee15@gmail.com",
	}

	err := s.userService.CreateUser(&user, nil)
	s.NoError(err)
	s.NotEmpty(user.ID)

	s.Run("create new user auth", func() {
		password, err := utils.GeneratePasswordHash("password")
		s.NoError(err)

		auth := models.Auth{
			UserID:   user.ID,
			Password: password,
		}

		then := time.Now()
		err = s.authService.CreateAuth(&auth, nil)
		s.NoError(err)

		s.NotEmpty(auth.ID)
		s.NotEmpty(auth.UpdatedAt)
		s.NotEmpty(auth.CreatedAt)
		s.Greater(auth.UpdatedAt, then)

		err = s.authService.CreateAuth(&auth, nil)
		s.Error(err)
		s.Contains(err.Error(), "duplicate key")
	})

	s.Run("get user auth", func() {
		auth, err := s.authService.GetUserAuth(user.ID, nil)
		s.NoError(err)

		s.Equal(user.ID, auth.UserID)
	})

	s.Run("update user auth", func() {
		password, err := utils.GeneratePasswordHash("passwordss")
		s.NoError(err)

		auth, err := s.authService.GetUserAuth(user.ID, nil)
		s.NoError(err)

		then := auth.UpdatedAt
		oldPassword := auth.Password
		auth.Password = password
		err = s.authService.UpdateUserAuth(auth, nil)
		s.NoError(err)

		s.NotEqual(password, oldPassword)
		s.Greater(auth.UpdatedAt, then)
	})
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
