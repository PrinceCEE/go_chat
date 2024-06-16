package websocket

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
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/princecee/go_chat/app/api/auth"
	"github.com/princecee/go_chat/app/api/rooms"
	"github.com/princecee/go_chat/app/api/users"
	"github.com/princecee/go_chat/internal/models"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
	"github.com/stretchr/testify/suite"
)

type WebsocketTestSuite struct {
	suite.Suite
	services         services.Services
	server           *httptest.Server
	sender, receiver struct {
		accessToken string
		user        *models.User
	}
}

func (s *WebsocketTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}

	conn, err := pgxpool.New(context.Background(), os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}

	s.services = services.New(conn)

	r := gin.New()

	auth.Routes(r.Group("/api/v1/auth"), s.services)
	rooms.Routes(r.Group("/api/v1/rooms"), s.services)
	users.Routes(r.Group("/api/v1/users"), s.services)

	SetupWebsocket(r, conn)

	s.server = httptest.NewServer(r.Handler())
}

func (s *WebsocketTestSuite) TearDownSuite() {
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

func (s *WebsocketTestSuite) TestWebsocket() {
	baseUrl := s.server.URL
	client := s.server.Client()

	signupDtos := []map[string]string{
		{
			"first_name": "Chimezie",
			"last_name":  "Edeh",
			"email":      "princecee15@gmail.com",
			"password":   "password",
		},
		{
			"first_name": "Yung",
			"last_name":  "Yu",
			"email":      "yungyu@gmail.com",
			"password":   "password",
		},
	}

	var roomID string
	for i, dto := range signupDtos {
		signupDtoJson, err := json.Marshal(dto)
		s.NoError(err)

		resp, err := client.Post(baseUrl+"/api/v1/auth/sign-up", "application/json", bytes.NewBuffer(signupDtoJson))
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(http.StatusOK, resp.StatusCode)

		user, accessToken, err := getUserDetails(resp)
		s.NoError(err)

		s.NotEmpty(&user)
		s.NotEmpty(&accessToken)
		s.Equal(dto["email"], user.Email)

		if i == 0 {
			s.sender.user = &user
			s.sender.accessToken = accessToken

			roomID, err = s.createRoom(baseUrl, client)
			s.NoError(err)
		} else if i == 1 {
			s.receiver.user = &user
			s.receiver.accessToken = accessToken

			err = s.joinRoom(baseUrl, roomID, client)
			s.NoError(err)
		}
	}

	s.Run("send and receive message in rooms", func() {
		url := fmt.Sprintf("ws://%s/ws", s.server.URL[7:]) // remove `http://` and replace with `ws://`
		fmt.Println(url)
		sheaders, rheaders := http.Header{}, http.Header{}
		sheaders.Add("Authorization", s.sender.accessToken)
		rheaders.Add("Authorization", s.receiver.accessToken)

		senderConn, _, err := websocket.DefaultDialer.Dial(url, sheaders)
		s.NoError(err)

		receiverConn, _, err := websocket.DefaultDialer.Dial(url, rheaders)
		s.NoError(err)

		defer senderConn.Close()
		defer receiverConn.Close()

		senderMsg := Message{
			RoomID:  roomID,
			UserID:  s.sender.user.ID,
			Content: "Hello boy! How are you",
		}
		err = senderConn.WriteJSON(&senderMsg)
		s.NoError(err)

		var receiverMsg Message
		err = receiverConn.ReadJSON(&receiverMsg)
		s.NoError(err)

		s.Equal(senderMsg.Content, receiverMsg.Content)
	})
}

func (s *WebsocketTestSuite) joinRoom(baseUrl, roomID string, client *http.Client) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/rooms/%s/join", baseUrl, roomID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", s.receiver.accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	s.Equal(http.StatusOK, resp.StatusCode)

	return nil
}

func (s *WebsocketTestSuite) createRoom(baseUrl string, client *http.Client) (string, error) {
	createRoomJson, err := json.Marshal(map[string]any{
		"name":        "Physics",
		"description": "all about physics",
		"max_members": 3,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/rooms", baseUrl), bytes.NewBuffer(createRoomJson))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", s.sender.accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var data utils.Response[map[string]models.Room]
	err = utils.ReadJSON(resp.Body, &data)
	if err != nil {
		return "", err
	}

	s.Equal(http.StatusOK, resp.StatusCode)

	return data.Data["room"].ID, nil
}

func getUserDetails(resp *http.Response) (models.User, string, error) {
	var data utils.Response[map[string]models.User]
	err := utils.ReadJSON(resp.Body, &data)
	if err != nil {
		return models.User{}, "", err
	}

	return data.Data["user"], data.Meta.AccessToken, nil
}

func TestWebsocket(t *testing.T) {
	suite.Run(t, new(WebsocketTestSuite))
}
