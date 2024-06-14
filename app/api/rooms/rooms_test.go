package rooms

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

type roomsTestSuite struct {
	suite.Suite
	services services.Services
	server   *httptest.Server
}

func (s *roomsTestSuite) SetupSuite() {
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
	Routes(r.Group("/api/v1/rooms"), s.services)

	s.server = httptest.NewServer(r.Handler())
}

func (s *roomsTestSuite) TearDownSuite() {
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

func (s *roomsTestSuite) TestRoomsHandler() {
	client := s.server.Client()
	baseUrl := s.server.URL
	roomBaseUrl := baseUrl + "/api/v1/rooms"
	contentType := "application/json"

	creatorDto := map[string]string{
		"first_name": "Chimezie",
		"last_name":  "Edeh",
		"email":      "princecee15@gmail.com",
		"password":   "password",
	}

	members := []*struct {
		signupDto   map[string]string
		user        models.User
		accessToken string
		statusCode  int
		message     string
		success     bool
	}{
		{
			signupDto: map[string]string{
				"first_name": "Miracle",
				"last_name":  "Ekene",
				"email":      "miracle@gmail.com",
				"password":   "password",
			},
			statusCode: http.StatusOK,
			success:    true,
			message:    "joined room successfully",
		},
		{
			signupDto: map[string]string{
				"first_name": "Yu",
				"last_name":  "Jong",
				"email":      "jong@gmail.com",
				"password":   "password",
			},
			statusCode: http.StatusUnauthorized,
			success:    false,
			message:    "max room members reached",
		},
	}

	creatorJson, err := json.Marshal(creatorDto)
	s.NoError(err)

	resp, err := client.Post(baseUrl+"/api/v1/auth/sign-up", contentType, bytes.NewBuffer(creatorJson))
	s.NoError(err)

	var data utils.Response[map[string]models.User]
	err = utils.ReadJSON(resp.Body, &data)
	s.NoError(err)

	user := data.Data["user"]
	accessToken := data.Meta.AccessToken

	s.NotEmpty(user)
	s.NotEmpty(accessToken)

	var room models.Room
	s.Run("create room", func() {
		createRoomDtoJson, err := json.Marshal(map[string]any{
			"name":        "Maths group",
			"description": "all about advanced maths",
			"max_members": 2,
		})

		s.NoError(err)

		req, err := http.NewRequest("POST", roomBaseUrl, bytes.NewBuffer(createRoomDtoJson))
		s.NoError(err)

		req.Header.Set("Authorization", accessToken)
		resp, err := client.Do(req)
		s.NoError(err)

		var data utils.Response[map[string]models.Room]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(http.StatusOK, resp.StatusCode)
		s.Equal("room created successfully", data.Message)
		s.Equal("Maths group", data.Data["room"].Name)

		room = data.Data["room"]
	})

	s.Run("get room", func() {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", roomBaseUrl, room.ID), nil)
		s.NoError(err)

		req.Header.Set("Authorization", accessToken)
		resp, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusOK, resp.StatusCode)

		var data utils.Response[map[string]models.Room]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(true, data.Success)
		s.Equal("fetched room successfully", data.Message)
		s.Equal(room.ID, data.Data["room"].ID)
	})

	s.Run("get rooms", func() {
		req, err := http.NewRequest("GET", roomBaseUrl, nil)
		s.NoError(err)

		req.Header.Set("Authorization", accessToken)
		resp, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusOK, resp.StatusCode)

		var data utils.Response[map[string][]models.Room]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(true, data.Success)
		s.Equal("fetched rooms successfully", data.Message)
		s.Equal(1, len(data.Data["rooms"]))
		s.Equal(room.ID, data.Data["rooms"][0].ID)
	})

	s.Run("update room", func() {
		updateDtoJson, err := json.Marshal(map[string]string{
			"name":        "Physics group",
			"description": "All about physics",
		})
		s.NoError(err)

		req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", roomBaseUrl, room.ID), bytes.NewBuffer(updateDtoJson))
		s.NoError(err)

		req.Header.Set("Authorization", accessToken)
		resp, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusOK, resp.StatusCode)
		var data utils.Response[map[string]models.Room]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)

		s.Equal("updated room successfully", data.Message)
		s.Equal(true, data.Success)
		s.Equal(room.ID, data.Data["room"].ID)
		s.Greater(data.Data["room"].UpdatedAt, room.UpdatedAt)
	})

	s.Run("join room", func() {
		for _, data := range members {
			signupDtoJson, err := json.Marshal(data.signupDto)
			s.NoError(err)

			resp, err := client.Post(baseUrl+"/api/v1/auth/sign-up", contentType, bytes.NewBuffer(signupDtoJson))
			s.NoError(err)

			s.Equal(http.StatusOK, resp.StatusCode)

			var _data utils.Response[map[string]models.User]
			err = utils.ReadJSON(resp.Body, &_data)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(true, _data.Success)
			data.user = _data.Data["user"]
			data.accessToken = _data.Meta.AccessToken
		}

		url := fmt.Sprintf("%s/%s/join", roomBaseUrl, room.ID)
		req, err := http.NewRequest("POST", url, nil)
		s.NoError(err)

		req.Header.Set("Authorization", accessToken)
		resp, err := client.Do(req)
		s.NoError(err)

		var data utils.Response[any]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(http.StatusNotAcceptable, resp.StatusCode)
		s.Equal("already a member", data.Message)

		for _, data := range members {
			url := fmt.Sprintf("%s/%s/join", roomBaseUrl, room.ID)
			req, err := http.NewRequest("POST", url, nil)
			s.NoError(err)

			req.Header.Set("Authorization", data.accessToken)
			resp, err := client.Do(req)
			s.NoError(err)

			var _data utils.Response[any]
			err = utils.ReadJSON(resp.Body, &_data)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(data.statusCode, resp.StatusCode)
			s.Equal(data.message, _data.Message)
			s.Equal(data.success, _data.Success)
		}
	})

	s.Run("leave room", func() {
		memberToLeave := members[0]

		url := fmt.Sprintf("%s/%s/leave", roomBaseUrl, room.ID)
		req, err := http.NewRequest("POST", url, nil)
		s.NoError(err)

		req.Header.Set("Authorization", memberToLeave.accessToken)
		resp, err := client.Do(req)
		s.NoError(err)

		var data utils.Response[any]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(http.StatusOK, resp.StatusCode)
		s.Equal("left room successfully", data.Message)
		s.Equal(true, data.Success)
	})

	s.Run("get room members", func() {
		url := fmt.Sprintf("%s/%s/members", roomBaseUrl, room.ID)
		req, err := http.NewRequest("GET", url, nil)
		s.NoError(err)

		req.Header.Set("Authorization", accessToken)
		resp, err := client.Do(req)
		s.NoError(err)

		var data utils.Response[map[string][]models.RoomMember]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(http.StatusOK, resp.StatusCode)
		s.Equal("members fetched successfully", data.Message)
		s.Equal(true, data.Success)
		s.Equal(1, len(data.Data["members"]))
		s.Equal(user.ID, data.Data["members"][0].UserID)
	})

	s.Run("get room messages", func() {
		url := fmt.Sprintf("%s/%s/messages", roomBaseUrl, room.ID)
		req, err := http.NewRequest("GET", url, nil)
		s.NoError(err)

		req.Header.Set("Authorization", accessToken)
		resp, err := client.Do(req)
		s.NoError(err)

		var data utils.Response[map[string][]models.RoomMessage]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(http.StatusOK, resp.StatusCode)
		s.Equal("messages fetched successfully", data.Message)
		s.Equal(true, data.Success)
		s.Equal(0, len(data.Data["messages"]))
	})

	s.Run("delete room", func() {
		url := fmt.Sprintf("%s/%s", roomBaseUrl, room.ID)
		req, err := http.NewRequest("DELETE", url, nil)
		s.NoError(err)

		req.Header.Set("Authorization", accessToken)
		resp, err := client.Do(req)
		s.NoError(err)

		var data utils.Response[any]
		err = utils.ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(http.StatusOK, resp.StatusCode)
		s.Equal("deleted room successfully", data.Message)
		s.Equal(true, data.Success)
	})
}

func TestRoomsHandler(t *testing.T) {
	suite.Run(t, new(roomsTestSuite))
}
