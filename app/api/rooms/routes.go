package rooms

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/middlewares"
	"github.com/princecee/go_chat/internal/services"
)

func Routes(r *gin.RouterGroup, s services.Services) {
	h := roomHandler{services: s}

	r.Use(middlewares.Authenticator(s))

	r.GET("/:roomId", middlewares.ErrorHandler(h.getRoom))
	r.GET("/", middlewares.ErrorHandler(h.getRooms))
	r.POST("/", middlewares.ErrorHandler(h.createRoom))
	r.PATCH("/:roomId", middlewares.ErrorHandler(h.updateRoom))
	r.DELETE("/:roomId", middlewares.ErrorHandler(h.deleteRoom))
	r.POST("/:roomId/join", middlewares.ErrorHandler(h.joinRoom))
	r.POST("/:roomId/leave", middlewares.ErrorHandler(h.leaveRoom))
	r.GET("/:roomId/members", middlewares.ErrorHandler(h.getRoomMembers))
	r.GET("/:roomId/messages", middlewares.ErrorHandler(h.getRoomMessages))
}
