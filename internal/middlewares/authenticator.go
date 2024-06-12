package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/internal/db/repositories"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
)

func Authenticator(s services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		arr := strings.Split(c.Request.Header.Get("Authorization"), " ")

		badReqResponse := utils.ResponseGeneric{
			Success: false,
			Message: "invalid token",
		}
		if len(arr) != 2 {
			c.JSON(http.StatusBadRequest, badReqResponse)
			return
		}

		token, err := utils.VerifyToken(arr[1])
		if err != nil {
			c.JSON(http.StatusBadRequest, badReqResponse)
			return
		}

		user, err := s.GetUserService().GetUser(repositories.GetUserParams{
			ID: token.ID,
		}, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, badReqResponse)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
