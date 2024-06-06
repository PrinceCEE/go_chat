package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
)

type authHandler struct {
	services services.Services
}

func (h *authHandler) signUp(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}

func (h *authHandler) signIn(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}

func (h *authHandler) resetPassword(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}

func (h *authHandler) changePassword(c *gin.Context) error {
	return &utils.ServerError{
		Message:    "not implemented",
		StatusCode: http.StatusNotImplemented,
	}
}
