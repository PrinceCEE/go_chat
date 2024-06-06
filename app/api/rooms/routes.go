package rooms

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/middlewares"
	"github.com/princecee/go_chat/internal/services"
)

func Routes(r *gin.RouterGroup, s services.Services) {
	h := roomHandler{services: s}

	r.GET("/:id", middlewares.ErrorHandler(h.getRoom))
	r.GET("/", middlewares.ErrorHandler(h.getRooms))
	r.POST("/", middlewares.ErrorHandler(h.createRoom))
	r.PATCH("/:id", middlewares.ErrorHandler(h.updateRoom))
	r.DELETE("/:id", middlewares.ErrorHandler(h.deleteRoom))
	r.POST("/:id/join", middlewares.ErrorHandler(h.joinRoom))
}
