package websocket

import "github.com/gin-gonic/gin"

func SetupWebsocket(r *gin.Engine) {
	wSocket := websocket{}

	r.GET("/ws", wSocket.handleHandshake)
}

type websocket struct{}

func (wSocket *websocket) handleHandshake(c *gin.Context) {}
