package auth

import "github.com/gin-gonic/gin"

type authHandler struct{}

func (h *authHandler) signUp(c *gin.Context) {}

func (h *authHandler) signIn(c *gin.Context) {}

func (h *authHandler) resetPassword(c *gin.Context) {}

func (h *authHandler) changePassword(c *gin.Context) {}
