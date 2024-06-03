package api

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/app/api/auth"
	"github.com/princecee/go_chat/app/api/rooms"
	"github.com/princecee/go_chat/app/api/users"
)

func SetupAPI(r *gin.Engine) {
	v1 := r.Group("/v1")

	auth.Routes(v1.Group("/auth"))
	rooms.Routes(v1.Group("/rooms"))
	users.Routes(v1.Group("/users"))
}
