package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/princecee/go_chat/internal/models"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	suite.Suite
	services services.Services
	server   *httptest.Server
}

func (s *AuthSuite) SetupSuite() {
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
	Routes(&r.RouterGroup, s.services)

	s.server = httptest.NewServer(r.Handler())
}

func (s *AuthSuite) TearDownSuite() {
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

func (s *AuthSuite) TestAuth() {
	client := s.server.Client()
	baseUrl := s.server.URL
	contentType := "application/json"

	user := map[string]string{
		"first_name": "Chimezie",
		"last_name":  "Edeh",
		"email":      "princecee15@gmail.com",
		"password":   "password",
	}

	s.Run("sign up", func() {
		userJson, err := json.Marshal(user)
		s.NoError(err)

		resp, err := client.Post(baseUrl+"/sign-up", contentType, bytes.NewBuffer(userJson))
		s.NoError(err)

		s.Equal(http.StatusOK, resp.StatusCode)

		var data utils.Response[map[string]models.User]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(true, data.Success)
		s.Equal("sign up successful", data.Message)
		s.Greater(len(data.Meta.AccessToken), 0)
		s.Equal(user["first_name"], data.Data["user"].FirstName)
	})

	s.Run("sign up already existing account", func() {
		userJson, err := json.Marshal(user)
		s.NoError(err)

		resp, err := client.Post(baseUrl+"/sign-up", contentType, bytes.NewBuffer(userJson))
		s.NoError(err)
		s.Equal(http.StatusBadRequest, resp.StatusCode)

		var data utils.Response[map[string]models.User]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.False(data.Success)
		s.Equal("account already exists", data.Message)
	})

	s.Run("sign in", func() {
		signInDto := map[string]string{"email": user["email"], "password": user["password"]}
		signInjson, err := json.Marshal(signInDto)
		s.NoError(err)

		resp, err := client.Post(baseUrl+"/sign-in", contentType, bytes.NewBuffer(signInjson))
		s.NoError(err)

		var data utils.Response[map[string]models.User]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(true, data.Success)
		s.Equal("sign in successful", data.Message)
		s.Equal(user["first_name"], data.Data["user"].FirstName)
		s.Greater(len(data.Meta.AccessToken), 1)
	})
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
