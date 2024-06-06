package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/go_chat/app/api/auth"
	"github.com/princecee/go_chat/app/api/rooms"
	"github.com/princecee/go_chat/app/api/users"
	"github.com/princecee/go_chat/internal/services"
)

func SetupAPI(r *gin.Engine, conn *pgxpool.Pool) {
	v1 := r.Group("/api/v1")

	services := services.New(conn)
	auth.Routes(v1.Group("/auth"), services)
	rooms.Routes(v1.Group("/rooms"), services)
	users.Routes(v1.Group("/users"), services)
}
