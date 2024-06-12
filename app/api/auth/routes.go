package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/middlewares"
	"github.com/princecee/go_chat/internal/services"
)

func Routes(r *gin.RouterGroup, s services.Services) {
	h := authHandler{services: s}

	r.POST("/sign-up", middlewares.ErrorHandler(h.signUp))
	r.POST("/sign-in", middlewares.ErrorHandler(h.signIn))
}
