package auth

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	h := authHandler{}

	r.POST("/sign-up", h.signUp)
	r.POST("/sign-in", h.signIn)
	r.POST("/reset-password", h.resetPassword)
	r.POST("/change-password", h.changePassword)
}
