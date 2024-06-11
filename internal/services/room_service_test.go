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
	"github.com/stretchr/testify/suite"
)

type RoomServiceTestSuite struct {
	suite.Suite
	conn        *pgxpool.Pool
	userService *userService
	roomService *roomService
}

func (s *RoomServiceTestSuite) SetupSuite() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}

	conn, err := pgxpool.New(context.Background(), os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}

	roomRepository := repositories.NewRoomRepository(conn)
	userRepository := repositories.NewUserRepository(conn)
	s.conn = conn
	s.roomService = &roomService{conn: conn, RoomRepository: roomRepository}
	s.userService = &userService{conn: conn, UserRepository: userRepository}
}

func (s *RoomServiceTestSuite) TearDownSuite() {
	defer s.conn.Close()

	teardownQuery := `
		DELETE FROM auths;
		DELETE FROM room_messages;
		DELETE FROM room_members;
		DELETE FROM rooms;
		DELETE FROM users;
	`

	_, err := s.conn.Exec(context.Background(), teardownQuery)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			panic(err)
		}
	}
}

func (s *RoomServiceTestSuite) TestRoomService() {
	users := []*models.User{
		{
			FirstName: "First 1",
			LastName:  "Last 1",
			Email:     "first1@example.com",
		},
		{
			FirstName: "First 2",
			LastName:  "Last 2",
			Email:     "first2@example.com",
		},
		{
			FirstName: "First 3",
			LastName:  "Last 3",
			Email:     "first3@example.com",
		},
		{
			FirstName: "First 4",
			LastName:  "Last 4",
			Email:     "first4@example.com",
		},
		{
			FirstName: "First 5",
			LastName:  "Last 5",
			Email:     "first5@example.com",
		},
	}

	for _, user := range users {
		err := s.userService.CreateUser(user, nil)
		s.NoError(err)
	}

	creator := users[0]
	var room *models.Room
	var roomMember *models.RoomMember

	s.Run("create room", func() {
		room = &models.Room{
			Name:        "Messi fans",
			Description: "Discussion about everything Messi",
			MaxMembers:  3,
			CreatedBy:   creator.ID,
		}

		err := s.roomService.CreateRoom(room, nil)
		s.NoError(err)

		s.NotEmpty(room.ID)
		s.NotEmpty(room.CreatedAt)
		s.NotEmpty(room.UpdatedAt)

		roomMemberCount, err := s.roomService.RoomRepository.GetRoomMemberCount(room.ID, nil)
		s.NoError(err)

		mc := *roomMemberCount
		s.Equal(1, mc)

		roomMembers, err := s.roomService.GetRoomMembers(repositories.GetRoomMembersParams{
			RoomID: &room.ID,
		}, nil)
		s.NoError(err)

		roomMember = roomMembers[0]
		s.Len(roomMembers, 1)
		s.Equal(room.ID, roomMember.RoomID)
		s.Equal(creator.ID, roomMember.UserID)
	})

	s.Run("get rooms", func() {
		rooms, err := s.roomService.GetRooms(nil, nil)
		s.NoError(err)
		s.Equal(room.ID, rooms[0].ID)
		s.Equal(room.CreatedBy, rooms[0].CreatedBy)
		s.Equal(1, len(rooms))

		rooms, err = s.roomService.GetRooms(&room.CreatedBy, nil)
		s.NoError(err)
		s.Equal(room.ID, rooms[0].ID)
		s.Equal(room.CreatedBy, rooms[0].CreatedBy)
		s.Equal(1, len(rooms))
	})

	s.Run("update room", func() {
		room.Description = "ordinary differential equations discussions"
		room.Name = "calculus group"

		then := room.UpdatedAt
		err := s.roomService.UpdateRoom(room, nil)
		s.NoError(err)
		s.Greater(room.UpdatedAt, then)
		s.Equal("calculus group", room.Name)
	})

	s.Run("join and leave room", func() {
		successUsers := users[1:3]
		errorUsers := users[3:]

		s.Run("join room", func() {
			for _, user := range successUsers {
				member := models.RoomMember{
					RoomID: room.ID,
					UserID: user.ID,
				}

				then := time.Now()
				err := s.roomService.JoinRoom(&member, nil)
				s.NoError(err)
				s.Greater(member.UpdatedAt, then)
			}

			membersCount, err := s.roomService.RoomRepository.GetRoomMemberCount(room.ID, nil)
			s.NoError(err)
			mc := *membersCount
			s.Equal(3, mc)

			for _, user := range errorUsers {
				member := models.RoomMember{
					RoomID: room.ID,
					UserID: user.ID,
				}

				err := s.roomService.JoinRoom(&member, nil)
				s.Error(err)
				s.ErrorIs(err, ErrMaxMembersReached)
			}
		})

		s.Run("leave room", func() {
			members, err := s.roomService.GetRoomMembers(repositories.GetRoomMembersParams{
				RoomID: &room.ID,
			}, nil)
			s.NoError(err)

			var member *models.RoomMember
			for _, v := range members {
				if v.UserID != creator.ID {
					member = v
					break
				}
			}

			err = s.roomService.LeaveRoom(member.ID, nil)
			s.NoError(err)

			membersCount, err := s.roomService.RoomRepository.GetRoomMemberCount(room.ID, nil)
			s.NoError(err)
			mc := *membersCount
			s.Equal(2, mc)

			newMember := models.RoomMember{
				RoomID: room.ID,
				UserID: errorUsers[0].ID,
			}

			err = s.roomService.JoinRoom(&newMember, nil)
			s.NoError(err)
		})
	})

	s.Run("manage messages", func() {
		message := models.RoomMessage{
			RoomID:       room.ID,
			RoomMemberID: roomMember.ID,
			UserID:       creator.ID,
			Content:      "Hello group",
		}

		s.Run("send message", func() {
			then := time.Now()
			err := s.roomService.CreateMessage(&message, nil)
			s.NoError(err)
			s.Greater(message.UpdatedAt, then)
			s.NotEmpty(message.ID)
		})

		s.Run("get message", func() {
			_message, err := s.roomService.GetMessage(message.ID, nil)
			s.NoError(err)
			s.Equal(message.ID, _message.ID)
			s.Equal(message.Content, _message.Content)
		})

		s.Run("get messages", func() {
			messages, err := s.roomService.GetMessages(repositories.GetRoomMessagesParams{
				RoomID: &room.ID,
				UserID: &creator.ID,
			}, nil)
			s.NoError(err)
			s.Equal(1, len(messages))
			s.Equal(message.ID, messages[0].ID)
			s.Equal(message.Content, messages[0].Content)
		})

		s.Run("delete message", func() {
			err := s.roomService.DeleteMessage(message.ID, nil)
			s.NoError(err)

			_message, err := s.roomService.GetMessage(message.ID, nil)
			s.ErrorIs(err, pgx.ErrNoRows)
			s.Empty(_message)
		})
	})

	s.Run("delete room", func() {
		err := s.roomService.DeleteRoom(room.ID, nil)
		s.NoError(err)

		_room, err := s.roomService.GetRoom(room.ID, nil)
		s.ErrorIs(err, pgx.ErrNoRows)
		s.Empty(_room)
	})
}

func TestRoomService(t *testing.T) {
	suite.Run(t, new(RoomServiceTestSuite))
}
