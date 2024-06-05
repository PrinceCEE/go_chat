package rooms

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/services"
)

type roomHandler struct {
	services services.Services
}

func (h *roomHandler) getRoom(c *gin.Context) {}

func (h *roomHandler) getRooms(c *gin.Context) {}

func (h *roomHandler) createRoom(c *gin.Context) {}

func (h *roomHandler) deleteRoom(c *gin.Context) {}

func (h *roomHandler) updateRoom(c *gin.Context) {}

func (h *roomHandler) joinRoom(c *gin.Context) {}
