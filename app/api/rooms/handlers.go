package rooms

import "github.com/gin-gonic/gin"

type roomHandler struct{}

func (h *roomHandler) getRoom(c *gin.Context) {}

func (h *roomHandler) getRooms(c *gin.Context) {}

func (h *roomHandler) createRoom(c *gin.Context) {}

func (h *roomHandler) deleteRoom(c *gin.Context) {}

func (h *roomHandler) updateRoom(c *gin.Context) {}

func (h *roomHandler) joinRoom(c *gin.Context) {}
