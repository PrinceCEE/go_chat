package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/db/repositories"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
)

func Authenticator(s services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		badReqResponse := utils.ResponseGeneric{
			Success: false,
			Message: "invalid token",
		}

		if token == "" {
			c.JSON(http.StatusBadRequest, badReqResponse)
			return
		}

		payload, err := utils.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, badReqResponse)
			return
		}

		user, err := s.GetUserService().GetUser(repositories.GetUserParams{
			ID: payload.Subject,
		}, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, badReqResponse)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
