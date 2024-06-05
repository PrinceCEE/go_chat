package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/services"
)

func Routes(r *gin.RouterGroup, s services.Services) {
	h := authHandler{services: s}

	r.POST("/sign-up", h.signUp)
	r.POST("/sign-in", h.signIn)
	r.POST("/reset-password", h.resetPassword)
	r.POST("/change-password", h.changePassword)
}
