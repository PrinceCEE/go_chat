package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type roomRepository struct {
	conn *pgxpool.Pool
}

func NewRoomRepository(conn *pgxpool.Pool) *roomRepository {
	return &roomRepository{conn}
}

func (r *roomRepository) CreateRoom()        {}
func (r *roomRepository) GetRoom()           {}
func (r *roomRepository) GetRooms()          {}
func (r *roomRepository) DeleteRoom()        {}
func (r *roomRepository) UpdateRoom()        {}
func (r *roomRepository) CreateRoomMember()  {}
func (r *roomRepository) GetRoomMember()     {}
func (r *roomRepository) GetRoomMembers()    {}
func (r *roomRepository) DeleteRoomMember()  {}
func (r *roomRepository) UpdateRoomMember()  {}
func (r *roomRepository) CreateRoomMessage() {}
func (r *roomRepository) GetRoomMessage()    {}
func (r *roomRepository) GetRoomMessages()   {}
func (r *roomRepository) DeleteRoomMessage() {}
