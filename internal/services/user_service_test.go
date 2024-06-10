package services

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/princecee/go_chat/internal/db/repositories"
	"github.com/princecee/go_chat/internal/models"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	service *userService
}

func (s *UserServiceTestSuite) SetupSuite() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}

	conn, err := pgxpool.New(context.Background(), os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}

	userRepository := repositories.NewUserRepository(conn)
	s.service = &userService{conn: conn, UserRepository: userRepository}
}

func (s *UserServiceTestSuite) TearDownSuite() {
	defer s.service.conn.Close()

	teardownQuery := `
		DELETE FROM users;
		DELETE FROM rooms;
	`

	_, err := s.service.conn.Exec(context.Background(), teardownQuery)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			panic(err)
		}
	}
}

func (s *UserServiceTestSuite) TestUserService() {
	var globalUser models.User

	s.Run("create user", func() {
		user := models.User{
			FirstName: "Chimezie",
			LastName:  "Edeh",
			Email:     "princecee15@gmail.com",
		}

		s.Empty(user.ID)
		s.Empty(user.CreatedAt)
		s.Empty(user.UpdatedAt)

		err := s.service.CreateUser(&user, nil)
		s.NoError(err)

		s.NotEmpty(user.ID)
		s.NotEmpty(user.CreatedAt)
		s.NotEmpty(user.UpdatedAt)

		globalUser = user
	})

	s.Run("get user", func() {
		s.Run("get user by ID", func() {
			user, err := s.service.GetUser(repositories.GetUserParams{
				ID: globalUser.ID,
			}, nil)

			s.NoError(err)
			s.Equal(globalUser.ID, user.ID)
		})

		s.Run("get user by Email", func() {
			user, err := s.service.GetUser(repositories.GetUserParams{
				Email: globalUser.Email,
			}, nil)

			s.NoError(err)
			s.Equal(globalUser.Email, user.Email)
		})
	})

	s.Run("get users", func() {
		users, err := s.service.GetUsers(nil)
		s.NoError(err)

		s.Equal(1, len(users))
		s.Equal(globalUser.ID, users[0].ID)
	})

	s.Run("update user", func() {
		user, err := s.service.GetUser(repositories.GetUserParams{
			ID: globalUser.ID,
		}, nil)

		s.NoError(err)

		user.FirstName = "Chimezie Update"
		user.LastName = "Edeh Update"

		err = s.service.UpdateUser(user, nil)
		s.NoError(err)

		_user, err := s.service.GetUser(repositories.GetUserParams{
			ID: globalUser.ID,
		}, nil)

		s.NoError(err)

		s.Equal(user.FirstName, _user.FirstName)
		s.Equal(user.LastName, _user.LastName)
		s.Greater(user.UpdatedAt, globalUser.UpdatedAt)
	})

	s.Run("delete user", func() {
		err := s.service.DeleteUser(globalUser.ID, nil)
		s.NoError(err)

		user, err := s.service.GetUser(repositories.GetUserParams{
			ID: globalUser.ID,
		}, nil)

		s.Empty(user)
		s.ErrorIs(err, pgx.ErrNoRows)
	})
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
