package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/go_chat/internal/db/repositories"
	"github.com/princecee/go_chat/internal/middlewares"
	"github.com/princecee/go_chat/internal/services"
	"github.com/princecee/go_chat/utils"
)

// accepts websocket connection from any origin
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsHandler struct {
	services services.Services
	clients  map[string]*wsClient
}

func SetupWebsocket(r *gin.Engine, db *pgxpool.Pool) {
	services := services.New(db)
	h := wsHandler{services: services}

	r.GET("/ws", middlewares.Authenticator(services), h.handleHandshake)
}

func (h *wsHandler) handleHandshake(c *gin.Context) {
	payload, err := utils.GetTokenFromRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ResponseGeneric{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	user, err := h.services.GetUserService().GetUser(repositories.GetUserParams{
		ID: payload.Subject,
	}, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseGeneric{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseGeneric{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if h.clients == nil {
		h.clients = map[string]*wsClient{}
	}

	client := &wsClient{
		conn:    conn,
		user:    user,
		handler: h,
		mu:      &sync.Mutex{},
	}
	h.clients[user.ID] = client

	errChan := make(chan error)
	go client.run(errChan)

	err = <-errChan
	if err != nil {
		if closeErr := conn.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}
}
