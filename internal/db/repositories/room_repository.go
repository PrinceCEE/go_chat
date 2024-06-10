package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	dataSource "github.com/princecee/go_chat/internal/db/data-source"
	"github.com/princecee/go_chat/internal/models"
	"github.com/princecee/go_chat/utils"
)

type roomRepository struct {
	conn *pgxpool.Pool
}

func NewRoomRepository(conn *pgxpool.Pool) *roomRepository {
	return &roomRepository{conn}
}

func (r *roomRepository) CreateRoom(room *models.Room, tx *pgxpool.Tx) error {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	_room, err := ds.CreateRoom(context.Background(), dataSource.CreateRoomParams{
		Name:        room.Name,
		Description: utils.StringToText(room.Description),
		MaxMembers:  int32(room.MaxMembers),
		CreatedBy:   utils.StringToUUID(room.CreatedBy),
	})
	if err != nil {
		return err
	}

	room.CreatedAt = _room.CreatedAt
	room.UpdatedAt = _room.UpdatedAt
	room.ID = _room.ID.String()
	return nil
}

func (r *roomRepository) GetRoom(id string, tx *pgxpool.Tx) (*models.Room, error) {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	_room, err := ds.GetRoom(context.Background(), utils.StringToUUID(id))
	if err != nil {
		return nil, err
	}

	return &models.Room{
		ID:          utils.UUIDToString(_room.ID),
		CreatedAt:   _room.CreatedAt,
		UpdatedAt:   _room.UpdatedAt,
		Name:        _room.Name,
		Description: _room.Description.String,
		CreatedBy:   utils.UUIDToString(_room.CreatedBy),
	}, nil
}

func (r *roomRepository) GetRooms(createdBy string, tx *pgxpool.Tx) ([]*models.Room, error) {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	cUUID := pgtype.UUID{}
	err := cUUID.Scan(createdBy)
	if err != nil {
		return nil, err
	}

	_rooms, err := ds.GetRooms(context.Background(), cUUID)
	if err != nil {
		return nil, err
	}

	rooms := []*models.Room{}
	for _, _room := range _rooms {
		r := &models.Room{
			ID:          utils.UUIDToString(_room.ID),
			CreatedAt:   _room.CreatedAt,
			UpdatedAt:   _room.UpdatedAt,
			Name:        _room.Name,
			Description: _room.Description.String,
			CreatedBy:   utils.UUIDToString(_room.CreatedBy),
		}
		rooms = append(rooms, r)
	}

	return rooms, nil
}

func (r *roomRepository) DeleteRoom(id string, tx *pgxpool.Tx) error {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	return ds.DeleteRoom(context.Background(), utils.StringToUUID(id))
}

func (r *roomRepository) UpdateRoom(room *models.Room, tx *pgxpool.Tx) error {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	room.UpdatedAt = time.Now()
	return ds.UpdateRoom(context.Background(), dataSource.UpdateRoomParams{
		UpdatedAt:   room.UpdatedAt,
		Name:        room.Name,
		Description: utils.StringToText(room.Description),
		MaxMembers:  int32(room.MaxMembers),
		ID:          utils.StringToUUID(room.ID),
	})
}

func (r *roomRepository) CreateRoomMember(member *models.RoomMember, tx *pgxpool.Tx) error {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	_member, err := ds.CreateRoomMember(context.Background(), dataSource.CreateRoomMemberParams{
		RoomID: utils.StringToUUID(member.RoomID),
		UserID: utils.StringToUUID(member.UserID),
	})
	if err != nil {
		return err
	}

	member.ID = utils.UUIDToString(_member.ID)
	member.CreatedAt = _member.CreatedAt
	member.UpdatedAt = _member.UpdatedAt

	return nil
}

func (r *roomRepository) GetRoomMember(id string, tx *pgxpool.Tx) (*models.RoomMember, error) {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	_member, err := ds.GetRoomMember(context.Background(), utils.StringToUUID(id))
	if err != nil {
		return nil, err
	}

	return &models.RoomMember{
		ID:        utils.UUIDToString(_member.ID),
		CreatedAt: _member.CreatedAt,
		UpdatedAt: _member.UpdatedAt,
		RoomID:    utils.UUIDToString(_member.RoomID),
		UserID:    utils.UUIDToString(_member.UserID),
	}, nil
}

type GetRoomMembersParams struct {
	RoomID *string
	UserID *string
}

func (r *roomRepository) GetRoomMembers(params GetRoomMembersParams, tx *pgxpool.Tx) ([]*models.RoomMember, error) {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	var userIDUUID, roomIDUUID pgtype.UUID
	err := userIDUUID.Scan(params.UserID)
	if err != nil {
		return nil, err
	}

	err = roomIDUUID.Scan(params.RoomID)
	if err != nil {
		return nil, err
	}

	_roomMembers, err := ds.GetRoomMembers(context.Background(), dataSource.GetRoomMembersParams{
		UserID: userIDUUID,
		RoomID: roomIDUUID,
	})
	if err != nil {
		return nil, err
	}

	roomMembers := []*models.RoomMember{}
	for _, member := range _roomMembers {
		roomMembers = append(roomMembers, &models.RoomMember{
			ID:        utils.UUIDToString(member.ID),
			CreatedAt: member.CreatedAt,
			UpdatedAt: member.UpdatedAt,
			RoomID:    utils.UUIDToString(member.RoomID),
			UserID:    utils.UUIDToString(member.UserID),
		})
	}

	return roomMembers, nil
}

func (r *roomRepository) DeleteRoomMember(id string, tx *pgxpool.Tx) error {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	return ds.DeleteRoomMember(context.Background(), utils.StringToUUID(id))
}

func (r *roomRepository) CreateRoomMessage(message *models.RoomMessage, tx *pgxpool.Tx) error {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	_message, err := ds.CreateRoomMessage(context.Background(), dataSource.CreateRoomMessageParams{
		RoomID:       utils.StringToUUID(message.RoomID),
		RoomMemberID: utils.StringToUUID(message.RoomMemberID),
		UserID:       utils.StringToUUID(message.UserID),
		Content:      message.Content,
	})
	if err != nil {
		return err
	}

	message.ID = utils.UUIDToString(_message.ID)
	message.CreatedAt = _message.CreatedAt
	message.UpdatedAt = _message.UpdatedAt

	return nil
}

func (r *roomRepository) GetRoomMessage(id string, tx *pgxpool.Tx) (*models.RoomMessage, error) {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	_message, err := ds.GetRoomMessage(context.Background(), utils.StringToUUID(id))
	if err != nil {
		return nil, err
	}

	return &models.RoomMessage{
		ID:           _message.ID.String(),
		CreatedAt:    _message.CreatedAt,
		UpdatedAt:    _message.UpdatedAt,
		RoomID:       _message.RoomID.String(),
		RoomMemberID: _message.RoomMemberID.String(),
		UserID:       _message.UserID.String(),
		Content:      _message.Content,
	}, nil
}

type GetRoomMessagesParams struct {
	RoomID       *string
	RoomMemberID *string
	UserID       *string
}

func (r *roomRepository) GetRoomMessages(params GetRoomMessagesParams, tx *pgxpool.Tx) ([]*models.RoomMessage, error) {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	var roomIDUUID, roomMemberIDUUID, userIDUUID pgtype.UUID
	err := roomIDUUID.Scan(params.RoomID)
	if err != nil {
		return nil, err
	}
	err = roomMemberIDUUID.Scan(params.RoomMemberID)
	if err != nil {
		return nil, err
	}
	err = userIDUUID.Scan(params.UserID)
	if err != nil {
		return nil, err
	}

	_messages, err := ds.GetRoomMessages(context.Background(), dataSource.GetRoomMessagesParams{
		RoomID:       roomIDUUID,
		RoomMemberID: roomMemberIDUUID,
		UserID:       userIDUUID,
	})
	if err != nil {
		return nil, err
	}

	messages := []*models.RoomMessage{}
	for _, message := range _messages {
		messages = append(messages, &models.RoomMessage{
			ID:           message.ID.String(),
			CreatedAt:    message.CreatedAt,
			UpdatedAt:    message.UpdatedAt,
			RoomID:       message.RoomID.String(),
			RoomMemberID: message.RoomMemberID.String(),
			UserID:       message.UserID.String(),
			Content:      message.Content,
		})
	}

	return messages, nil
}

func (r *roomRepository) DeleteRoomMessage(id string, tx *pgxpool.Tx) error {
	ds := dataSource.New(r.conn)
	if tx != nil {
		ds = ds.WithTx(tx)
	}

	return ds.DeleteRoomMessage(context.Background(), utils.StringToUUID(id))
}
