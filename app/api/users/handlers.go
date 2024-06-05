package users

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/services"
)

type userHandler struct {
	services services.Services
}

func (h *userHandler) getAccount(c *gin.Context) {}

func (h *userHandler) updateAccount(c *gin.Context) {}

func (h *userHandler) deleteAccount(c *gin.Context) {}
