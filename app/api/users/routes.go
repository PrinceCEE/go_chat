package users

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/services"
)

func Routes(r *gin.RouterGroup, s services.Services) {
	h := userHandler{services: s}

	r.GET("/:id", h.getAccount)
	r.DELETE("/:id", h.deleteAccount)
	r.PATCH("/:id", h.updateAccount)
}
