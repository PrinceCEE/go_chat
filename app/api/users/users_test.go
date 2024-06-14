package users

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/princecee/go_chat/app/api/auth"
	"github.com/princecee/go_chat/internal/models"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
	"github.com/stretchr/testify/suite"
)

type usersTestSuite struct {
	suite.Suite
	services services.Services
	server   *httptest.Server
}

func (s *usersTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal(err)
	}

	conn, err := pgxpool.New(context.Background(), os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}

	s.services = services.New(conn)

	r := gin.New()

	auth.Routes(r.Group("/api/v1/auth"), s.services)
	Routes(r.Group("/api/v1/users"), s.services)

	s.server = httptest.NewServer(r.Handler())
}

func (s *usersTestSuite) TearDownSuite() {
	db := s.services.GetDB()
	defer db.Close()
	defer s.server.Close()

	teardownQuery := `
		DELETE FROM auths;
		DELETE FROM room_messages;
		DELETE FROM room_members;
		DELETE FROM rooms;
		DELETE FROM users;
	`

	_, err := db.Exec(context.Background(), teardownQuery)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			panic(err)
		}
	}
}

func (s *usersTestSuite) TestUsersHandler() {
	client := s.server.Client()
	baseUrl := s.server.URL
	contentType := "application/json"

	signUpDto := map[string]string{
		"first_name": "Chimezie",
		"last_name":  "Edeh",
		"email":      "princecee15@gmail.com",
		"password":   "password",
	}

	signupJson, err := json.Marshal(signUpDto)
	s.NoError(err)

	resp, err := client.Post(baseUrl+"/api/v1/auth/sign-up", contentType, bytes.NewBuffer(signupJson))
	s.NoError(err)

	var data utils.Response[map[string]models.User]
	err = utils.ReadJSON(resp.Body, &data)
	s.NoError(err)

	user := data.Data["user"]
	accessToken := data.Meta.AccessToken

	s.NotEmpty(user)
	s.NotEmpty(accessToken)

	s.Run("get user without access token", func() {
		url := fmt.Sprintf("%s/api/v1/users/%s", baseUrl, user.ID)

		resp, err := client.Get(url)
		s.NoError(err)

		s.Equal(http.StatusBadRequest, resp.StatusCode)

		var data utils.Response[struct {
			Success bool
			Message string
		}]

		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(false, data.Success)
		s.Equal("invalid token", data.Message)
	})

	s.Run("get user with bad access token", func() {
		url := fmt.Sprintf("%s/api/v1/users/%s", baseUrl, user.ID)
		req, err := http.NewRequest("GET", url, nil)
		s.NoError(err)

		req.Header.Set("Authorization", "asnaksdnkans")

		resp, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusBadRequest, resp.StatusCode)

		var data utils.Response[struct {
			Success bool
			Message string
		}]

		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(false, data.Success)
		s.Equal("invalid token", data.Message)
	})

	s.Run("get user with access token", func() {
		url := fmt.Sprintf("%s/api/v1/users/%s", baseUrl, user.ID)
		req, err := http.NewRequest("GET", url, nil)
		s.NoError(err)

		req.Header.Set("Authorization", accessToken)

		resp, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusOK, resp.StatusCode)

		var data utils.Response[map[string]models.User]

		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(true, data.Success)
		s.Equal("fetched user successfully", data.Message)
		s.Equal(signUpDto["email"], data.Data["user"].Email)
	})

	s.Run("update user", func() {
		url := fmt.Sprintf("%s/api/v1/users/%s", baseUrl, user.ID)

		updateAccountJson, err := json.Marshal(map[string]string{
			"first_name": "Chimezie Update",
			"last_name":  "Edeh Update",
			"email":      "prince@gmail.com",
		})
		s.NoError(err)

		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(updateAccountJson))
		s.NoError(err)
		req.Header.Set("Authorization", accessToken)
		resp, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusOK, resp.StatusCode)

		var data utils.Response[map[string]models.User]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal("user updated successfully", data.Message)
		s.NotEqual(user.FirstName, data.Data["user"].FirstName)
		s.Equal("Chimezie Update", data.Data["user"].FirstName)
	})

	s.Run("delete user", func() {
		url := fmt.Sprintf("%s/api/v1/users/%s", baseUrl, user.ID)

		req, err := http.NewRequest("DELETE", url, nil)
		s.NoError(err)
		req.Header.Set("Authorization", accessToken)
		resp, err := client.Do(req)
		s.NoError(err)

		var data utils.ResponseGeneric

		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(true, data.Success)
		s.Equal("user deleted successfully", data.Message)

		req, err = http.NewRequest("GET", url, nil)
		s.NoError(err)

		req.Header.Set("Authorization", accessToken)

		resp, err = client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusBadRequest, resp.StatusCode)

		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(false, data.Success)
		s.Equal("invalid token", data.Message)
	})
}

func TestUsersHandler(t *testing.T) {
	suite.Run(t, new(usersTestSuite))
}
