package rooms

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/services"
)

func Routes(r *gin.RouterGroup, s services.Services) {
	h := roomHandler{services: s}

	r.GET("/:id", h.getRoom)
	r.GET("/", h.getRooms)
	r.POST("/", h.createRoom)
	r.PATCH("/:id", h.updateRoom)
	r.DELETE("/:id", h.deleteRoom)
	r.POST("/:id/join", h.joinRoom)
}
