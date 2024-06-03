package rooms

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	h := roomHandler{}

	r.GET("/:id", h.getRoom)
	r.GET("/", h.getRooms)
	r.POST("/", h.createRoom)
	r.PATCH("/:id", h.updateRoom)
	r.DELETE("/:id", h.deleteRoom)
	r.POST("/:id/join", h.joinRoom)
}
