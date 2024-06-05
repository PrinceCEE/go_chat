package websocket

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

// accepts websocket connection from any origin
// hence no need adding function handler for
// checking the origin
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SetupWebsocket(r *gin.Engine, conn *pgxpool.Pool) {
	wSocket := websocketHandler{}

	r.GET("/ws", wSocket.handleHandshake)
}

type websocketHandler struct{}

func (h *websocketHandler) handleHandshake(c *gin.Context) {
	_, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO - set up websocket clients map
}
