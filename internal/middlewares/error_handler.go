package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/utils"
)

type HandlerFunc func(c *gin.Context) error

func ErrorHandler(fn HandlerFunc) func(*gin.Context) {
	return func(c *gin.Context) {
		err := fn(c).(*utils.ServerError)
		if err != nil {
			errResponse := utils.Response[any]{
				Message: err.Error(),
				Success: false,
			}
			c.JSON(err.StatusCode, errResponse)
		}
	}
}
