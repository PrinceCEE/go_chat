package services

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/go_chat/internal/db/repositories"
	"github.com/princecee/go_chat/internal/models"
)

var (
	ErrMaxMembersReached = errors.New("max room members reached")
)

type roomService struct {
	conn           *pgxpool.Pool
	RoomRepository RoomRepository
}

func NewRoomService(conn *pgxpool.Pool) RoomService {
	return &roomService{
		conn:           conn,
		RoomRepository: repositories.NewRoomRepository(conn),
	}
}

func (s *roomService) CreateRoom(room *models.Room, tx *pgxpool.Tx) error {
	err := s.RoomRepository.CreateRoom(room, tx)
	if err != nil {
		return err
	}

	member := &models.RoomMember{
		RoomID: room.ID,
		UserID: room.CreatedBy,
	}
	return s.JoinRoom(member, tx)
}

func (s *roomService) GetRoom(id string, tx *pgxpool.Tx) (*models.Room, error) {
	return s.RoomRepository.GetRoom(id, tx)
}

func (s *roomService) GetRooms(createdBy *string, tx *pgxpool.Tx) ([]*models.Room, error) {
	return s.RoomRepository.GetRooms(createdBy, tx)
}

func (s *roomService) DeleteRoom(id string, tx *pgxpool.Tx) error {
	return s.RoomRepository.DeleteRoom(id, tx)
}

func (s *roomService) UpdateRoom(room *models.Room, tx *pgxpool.Tx) error {
	return s.RoomRepository.UpdateRoom(room, tx)
}

func (s *roomService) GetRoomMember(id string, tx *pgxpool.Tx) (*models.RoomMember, error) {
	return s.RoomRepository.GetRoomMember(id, tx)
}

func (s *roomService) LeaveRoom(roomMemberID string, tx *pgxpool.Tx) error {
	return s.RoomRepository.DeleteRoomMember(roomMemberID, tx)
}

func (s *roomService) JoinRoom(member *models.RoomMember, tx *pgxpool.Tx) error {
	room, err := s.GetRoom(member.RoomID, tx)
	if err != nil {
		return err
	}

	memberCount, err := s.RoomRepository.GetRoomMemberCount(room.ID, tx)
	if err != nil {
		return err
	}

	if room.MaxMembers < *memberCount+1 {
		return ErrMaxMembersReached
	}

	return s.RoomRepository.CreateRoomMember(member, tx)
}

func (s *roomService) GetRoomMembers(params repositories.GetRoomMembersParams, tx *pgxpool.Tx) ([]*models.RoomMember, error) {
	return s.RoomRepository.GetRoomMembers(params, tx)
}

func (s *roomService) CreateMessage(message *models.RoomMessage, tx *pgxpool.Tx) error {
	return s.RoomRepository.CreateRoomMessage(message, tx)
}

func (s *roomService) GetMessage(id string, tx *pgxpool.Tx) (*models.RoomMessage, error) {
	return s.RoomRepository.GetRoomMessage(id, tx)
}

func (s *roomService) GetMessages(params repositories.GetRoomMessagesParams, tx *pgxpool.Tx) ([]*models.RoomMessage, error) {
	return s.RoomRepository.GetRoomMessages(params, tx)
}

func (s *roomService) DeleteMessage(id string, tx *pgxpool.Tx) error {
	return s.RoomRepository.DeleteRoomMessage(id, tx)
}

type RoomRepository interface {
	CreateRoom(room *models.Room, tx *pgxpool.Tx) error
	GetRoom(id string, tx *pgxpool.Tx) (*models.Room, error)
	GetRooms(createdBy *string, tx *pgxpool.Tx) ([]*models.Room, error)
	DeleteRoom(id string, tx *pgxpool.Tx) error
	UpdateRoom(room *models.Room, tx *pgxpool.Tx) error
	CreateRoomMember(member *models.RoomMember, tx *pgxpool.Tx) error
	GetRoomMember(id string, tx *pgxpool.Tx) (*models.RoomMember, error)
	GetRoomMemberCount(roomId string, tx *pgxpool.Tx) (*int, error)
	GetRoomMembers(params repositories.GetRoomMembersParams, tx *pgxpool.Tx) ([]*models.RoomMember, error)
	DeleteRoomMember(id string, tx *pgxpool.Tx) error
	CreateRoomMessage(message *models.RoomMessage, tx *pgxpool.Tx) error
	GetRoomMessage(id string, tx *pgxpool.Tx) (*models.RoomMessage, error)
	GetRoomMessages(params repositories.GetRoomMessagesParams, tx *pgxpool.Tx) ([]*models.RoomMessage, error)
	DeleteRoomMessage(id string, tx *pgxpool.Tx) error
}

type RoomService interface {
	CreateRoom(room *models.Room, tx *pgxpool.Tx) error
	GetRoom(id string, tx *pgxpool.Tx) (*models.Room, error)
	GetRooms(createdBy *string, tx *pgxpool.Tx) ([]*models.Room, error)
	DeleteRoom(id string, tx *pgxpool.Tx) error
	UpdateRoom(room *models.Room, tx *pgxpool.Tx) error
	GetRoomMember(id string, tx *pgxpool.Tx) (*models.RoomMember, error)
	LeaveRoom(id string, tx *pgxpool.Tx) error
	JoinRoom(member *models.RoomMember, tx *pgxpool.Tx) error
	GetRoomMembers(params repositories.GetRoomMembersParams, tx *pgxpool.Tx) ([]*models.RoomMember, error)
	CreateMessage(message *models.RoomMessage, tx *pgxpool.Tx) error
	GetMessage(id string, tx *pgxpool.Tx) (*models.RoomMessage, error)
	GetMessages(params repositories.GetRoomMessagesParams, tx *pgxpool.Tx) ([]*models.RoomMessage, error)
	DeleteMessage(id string, tx *pgxpool.Tx) error
}
