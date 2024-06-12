package rooms

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/princecee/go_chat/internal/db/repositories"
	"github.com/princecee/go_chat/internal/models"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
)

type roomHandler struct {
	services services.Services
}

func (h *roomHandler) getRoom(c *gin.Context) error {
	roomId := c.Params.ByName("roomId")
	val, ok := c.Get("user")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	user := val.(*models.User)

	roomService := h.services.GetRoomService()
	_, err := roomService.GetRoomMemberByWhere(repositories.GetRoomMemberByWhereParams{
		UserID: user.ID,
		RoomID: roomId,
	}, nil)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    utils.ErrUnauthorized.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	room, err := roomService.GetRoom(roomId, nil)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Success: true,
		Message: "fetched room successfully",
		Data:    map[string]any{"room": room},
	})
	return nil
}

func (h *roomHandler) getRooms(c *gin.Context) error {
	rooms, err := h.services.GetRoomService().GetRooms(nil, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusOK, utils.ResponseGeneric{
				Success: true,
				Message: "fetched rooms successfully",
				Data:    map[string]any{"rooms": []models.Room{}},
			})
			return nil
		}

		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Success: true,
		Message: "fetched rooms successfully",
		Data:    map[string]any{"rooms": rooms},
	})

	return nil
}

type CreateRoomDto struct {
	Name        string `json:"name" validate:"required,alphanumeric"`
	Description string `json:"description" validate:"required,alphanumeric"`
	MaxMembers  int    `json:"max_members" validate:"required,gt=0"`
}

func (h *roomHandler) createRoom(c *gin.Context) error {
	val, ok := c.Get("user")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	var createRoomDto CreateRoomDto
	err := c.Bind(&createRoomDto)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	user := val.(*models.User)
	room := &models.Room{
		CreatedBy:   user.ID,
		Description: createRoomDto.Description,
		Name:        createRoomDto.Name,
		MaxMembers:  createRoomDto.MaxMembers,
	}

	tx, _ := h.services.GetDB().Begin(context.Background())
	err = h.services.GetRoomService().CreateRoom(room, tx)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Success: true,
		Message: "room created successfully",
		Data:    map[string]*models.Room{"room": room},
	})

	return nil
}

func (h *roomHandler) deleteRoom(c *gin.Context) error {
	roomId := c.Params.ByName("roomId")

	val, ok := c.Get("user")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	user := val.(*models.User)

	roomService := h.services.GetRoomService()
	_, err := roomService.GetRoomMemberByWhere(repositories.GetRoomMemberByWhereParams{
		UserID: user.ID,
		RoomID: roomId,
	}, nil)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    utils.ErrUnauthorized.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	room, err := roomService.GetRoom(roomId, nil)
	if err != nil {
		return &utils.ServerError{
			Message:    err.Error(),
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if room.CreatedBy != user.ID {
		return &utils.ServerError{
			Message:    utils.ErrUnauthorized.Error(),
			Err:        utils.ErrUnauthorized,
			StatusCode: http.StatusUnauthorized,
		}
	}

	err = roomService.DeleteMessage(roomId, nil)
	if err != nil {
		return &utils.ServerError{
			Message:    utils.ErrUnauthorized.Error(),
			Err:        utils.ErrUnauthorized,
			StatusCode: http.StatusUnauthorized,
		}
	}

	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Success: true,
		Message: "deleted room successfully",
	})

	return nil
}

type UpdateRoomDto struct {
	Name        *string `json:"name,omitempty" validate:"alphanumeric"`
	Description *string `json:"description,omitempty" validate:"alphanumeric"`
	MaxMembers  *int    `json:"max_members,omitempty" validate:"gt=0"`
}

func (h *roomHandler) updateRoom(c *gin.Context) error {
	roomId := c.Params.ByName("roomId")

	val, ok := c.Get("user")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	var updateRoomDto UpdateRoomDto
	err := c.BindJSON(&updateRoomDto)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	user := val.(*models.User)

	roomService := h.services.GetRoomService()
	_, err = roomService.GetRoomMemberByWhere(repositories.GetRoomMemberByWhereParams{
		UserID: user.ID,
		RoomID: roomId,
	}, nil)
	if err != nil {
		return &utils.ServerError{
			Err:        err,
			Message:    utils.ErrUnauthorized.Error(),
			StatusCode: http.StatusUnauthorized,
		}
	}

	room, err := roomService.GetRoom(roomId, nil)
	if err != nil {
		return &utils.ServerError{
			Message:    err.Error(),
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if room.CreatedBy != user.ID {
		return &utils.ServerError{
			Message:    utils.ErrUnauthorized.Error(),
			Err:        utils.ErrUnauthorized,
			StatusCode: http.StatusUnauthorized,
		}
	}

	if updateRoomDto.Name != nil {
		room.Name = *updateRoomDto.Name
	}
	if updateRoomDto.Description != nil {
		room.Description = *updateRoomDto.Description
	}
	if updateRoomDto.MaxMembers != nil {
		room.MaxMembers = *updateRoomDto.MaxMembers
	}

	return nil
}

func (h *roomHandler) joinRoom(c *gin.Context) error {
	roomId := c.Params.ByName("roomId")

	val, ok := c.Get("user")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	user := val.(*models.User)
	roomService := h.services.GetRoomService()

	_, err := roomService.GetRoom(roomId, nil)
	if err != nil {
		se := utils.ServerError{Err: err}
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			se.Message = utils.ErrNotFound.Error()
			se.StatusCode = http.StatusNotFound
		default:
			se.Message = err.Error()
			se.StatusCode = http.StatusInternalServerError
		}

		return &se
	}

	member, err := roomService.GetRoomMemberByWhere(repositories.GetRoomMemberByWhereParams{
		UserID: user.ID,
		RoomID: roomId,
	}, nil)
	if err != nil || member != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return &utils.ServerError{
				Err:        err,
				Message:    utils.ErrUnauthorized.Error(),
				StatusCode: http.StatusUnauthorized,
			}
		}

		if member != nil {
			return &utils.ServerError{
				Err:        utils.ErrDuplicateRecord,
				Message:    "already a member",
				StatusCode: http.StatusNotAcceptable,
			}
		}
	}

	member = &models.RoomMember{
		UserID: user.ID,
		RoomID: roomId,
	}
	err = roomService.JoinRoom(member, nil)
	if err != nil {
		return &utils.ServerError{
			Message:    err.Error(),
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Success: true,
		Message: "joined room successfully",
		Data:    map[string]*models.RoomMember{"member": member},
	})
	return nil
}

func (h *roomHandler) leaveRoom(c *gin.Context) error {
	roomId := c.Params.ByName("roomId")

	val, ok := c.Get("user")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	user := val.(*models.User)
	roomService := h.services.GetRoomService()

	_, err := roomService.GetRoom(roomId, nil)
	if err != nil {
		se := utils.ServerError{Err: err}
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			se.Message = utils.ErrNotFound.Error()
			se.StatusCode = http.StatusNotFound
		default:
			se.Message = err.Error()
			se.StatusCode = http.StatusInternalServerError
		}

		return &se
	}

	member, err := roomService.GetRoomMemberByWhere(repositories.GetRoomMemberByWhereParams{
		UserID: user.ID,
		RoomID: roomId,
	}, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &utils.ServerError{
				Err:        err,
				Message:    "not a member of room",
				StatusCode: http.StatusBadRequest,
			}
		}

		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	err = roomService.LeaveRoom(member.ID, nil)
	if err != nil {
		return &utils.ServerError{
			Message:    err.Error(),
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Success: true,
		Message: "left room successfully",
	})
	return nil
}

func (h *roomHandler) getRoomMembers(c *gin.Context) error {
	roomId := c.Params.ByName("roomId")

	val, ok := c.Get("user")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	user := val.(*models.User)
	roomService := h.services.GetRoomService()

	_, err := roomService.GetRoom(roomId, nil)
	if err != nil {
		se := utils.ServerError{Err: err}
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			se.Message = utils.ErrNotFound.Error()
			se.StatusCode = http.StatusNotFound
		default:
			se.Message = err.Error()
			se.StatusCode = http.StatusInternalServerError
		}

		return &se
	}

	_, err = roomService.GetRoomMemberByWhere(repositories.GetRoomMemberByWhereParams{
		UserID: user.ID,
		RoomID: roomId,
	}, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &utils.ServerError{
				Err:        err,
				Message:    "not a member of room",
				StatusCode: http.StatusBadRequest,
			}
		}

		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	members, err := roomService.GetRoomMembers(repositories.GetRoomMembersParams{
		UserID: &user.ID,
		RoomID: &roomId,
	}, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusOK, utils.ResponseGeneric{
				Success: true,
				Message: "members fetched successfully",
				Data:    map[string][]*models.RoomMember{},
			})
			return nil
		}

		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Success: true,
		Message: "members fetched successfully",
		Data:    members,
	})

	return nil
}

func (h *roomHandler) getRoomMessages(c *gin.Context) error {
	roomId := c.Params.ByName("roomId")

	val, ok := c.Get("user")
	if !ok {
		return &utils.ServerError{
			Message:    utils.ErrInternalServer.Error(),
			Err:        utils.ErrInternalServer,
			StatusCode: http.StatusInternalServerError,
		}
	}

	user := val.(*models.User)
	roomService := h.services.GetRoomService()

	_, err := roomService.GetRoom(roomId, nil)
	if err != nil {
		se := utils.ServerError{Err: err}
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			se.Message = utils.ErrNotFound.Error()
			se.StatusCode = http.StatusNotFound
		default:
			se.Message = err.Error()
			se.StatusCode = http.StatusInternalServerError
		}

		return &se
	}

	_, err = roomService.GetRoomMemberByWhere(repositories.GetRoomMemberByWhereParams{
		UserID: user.ID,
		RoomID: roomId,
	}, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &utils.ServerError{
				Err:        err,
				Message:    "not a member of room",
				StatusCode: http.StatusBadRequest,
			}
		}

		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	messages, err := roomService.GetMessages(repositories.GetRoomMessagesParams{
		RoomID: &roomId,
	}, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusOK, utils.ResponseGeneric{
				Success: true,
				Message: "messages fetched successfully",
				Data:    map[string][]*models.RoomMessage{},
			})
			return nil
		}

		return &utils.ServerError{
			Err:        err,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	c.JSON(http.StatusOK, utils.ResponseGeneric{
		Success: true,
		Message: "messages fetched successfully",
		Data:    messages,
	})

	return nil
}
