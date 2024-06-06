package rooms

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
)

type roomHandler struct {
	services services.Services
}

func (h *roomHandler) getRoom(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}

func (h *roomHandler) getRooms(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}

func (h *roomHandler) createRoom(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}

func (h *roomHandler) deleteRoom(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}

func (h *roomHandler) updateRoom(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}

func (h *roomHandler) joinRoom(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}
