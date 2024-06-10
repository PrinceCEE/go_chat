package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/princecee/go_chat/utils"
	"github.com/stretchr/testify/assert"
)

func TestErrorhandler(t *testing.T) {
	errStr := "internal server error"

	fn1 := func(c *gin.Context) error {
		c.JSON(http.StatusOK, utils.Response[any]{
			Message: "success",
			Success: true,
		})
		return nil
	}

	fn2 := func(c *gin.Context) error {
		return &utils.ServerError{
			Message:    errStr,
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New(errStr),
		}
	}

	gin.SetMode(gin.TestMode)
	g := gin.New()

	g.GET("/success", ErrorHandler(fn1))
	g.GET("/error", ErrorHandler(fn2))

	server := httptest.NewServer(g.Handler())
	defer server.Close()

	client := server.Client()

	tests := []struct {
		fn                 func(string) (*http.Response, error)
		name               string
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
		url                string
	}{
		{
			fn: func(url string) (*http.Response, error) {
				return client.Get(url)
			},
			name:               "Test with error",
			expectedStatusCode: 500,
			expectedMessage:    errStr,
			expectedSuccess:    false,
			url:                fmt.Sprintf(server.URL + "/error"),
		},
		{
			fn: func(url string) (*http.Response, error) {
				return client.Get(url)
			},
			name:               "Test without error",
			expectedStatusCode: 200,
			expectedMessage:    "success",
			expectedSuccess:    true,
			url:                fmt.Sprintf(server.URL + "/success"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := test.fn(test.url)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, test.expectedStatusCode, resp.StatusCode)
			var data utils.Response[any]
			err = json.NewDecoder(resp.Body).Decode(&data)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedMessage, data.Message)
			assert.Equal(t, test.expectedSuccess, data.Success)
		})
	}
}
