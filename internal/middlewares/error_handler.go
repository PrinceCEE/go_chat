package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/utils"
)

type HandlerFunc func(c *gin.Context) error

func ErrorHandler(fn HandlerFunc) func(*gin.Context) {
	return func(c *gin.Context) {
		err := fn(c)

		if err != nil {
			_err, ok := err.(*utils.ServerError)
			errResponse := utils.Response[any]{
				Success: false,
				Message: err.Error(),
			}

			var statusCode int
			if !ok {
				statusCode = http.StatusInternalServerError
			} else {
				statusCode = _err.StatusCode
			}

			c.JSON(statusCode, errResponse)
		}
	}
}
