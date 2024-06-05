package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/services"
)

type authHandler struct {
	services services.Services
}

func (h *authHandler) signUp(c *gin.Context) {}

func (h *authHandler) signIn(c *gin.Context) {}

func (h *authHandler) resetPassword(c *gin.Context) {}

func (h *authHandler) changePassword(c *gin.Context) {}
