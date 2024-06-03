package users

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	h := userHandler{}

	r.GET("/:id", h.getAccount)
	r.DELETE("/:id", h.deleteAccount)
	r.PATCH("/:id", h.updateAccount)
}
