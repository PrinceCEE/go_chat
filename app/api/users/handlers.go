package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
)

type userHandler struct {
	services services.Services
}

func (h *userHandler) getAccount(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}

func (h *userHandler) updateAccount(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}

func (h *userHandler) deleteAccount(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}
