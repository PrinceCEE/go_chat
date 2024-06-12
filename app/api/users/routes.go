package users

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/middlewares"
	"github.com/princecee/go_chat/internal/services"
)

func Routes(r *gin.RouterGroup, s services.Services) {
	h := userHandler{services: s}

	r.Use(middlewares.Authenticator(s))

	r.GET("/:id", middlewares.ErrorHandler(h.getAccount))
	r.DELETE("/:id", middlewares.ErrorHandler(h.deleteAccount))
	r.PATCH("/:id", middlewares.ErrorHandler(h.updateAccount))
}
