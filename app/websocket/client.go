package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/princecee/go_chat/internal/db/repositories"
	"github.com/princecee/go_chat/internal/models"
)

type wsClient struct {
	conn    *websocket.Conn
	user    *models.User
	handler *wsHandler
	mu      *sync.Mutex
}

type Message struct {
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
	UserID  string `json:"user_id,omitempty"`
}

func (client *wsClient) run(errChan chan<- error) {
	defer close(errChan)

	for {
		data := new(Message)
		err := client.conn.ReadJSON(data)
		if err != nil {
			errChan <- err
			break
		}

		roomService := client.handler.services.GetRoomService()
		roomMember, err := roomService.GetRoomMemberByWhere(repositories.GetRoomMemberByWhereParams{
			UserID: client.user.ID,
			RoomID: data.RoomID,
		}, nil)
		if err != nil {
			errChan <- err
			break
		}

		message := &models.RoomMessage{
			RoomID:       data.RoomID,
			UserID:       client.user.ID,
			RoomMemberID: roomMember.ID,
			Content:      data.Content,
		}
		err = roomService.CreateMessage(message, nil)
		if err != nil {
			errChan <- err
			break
		}

		client.mu.Lock()
		err = client.broadcast(client.user, data.RoomID, data.Content)
		if err != nil {
			errChan <- err
			break
		}
		client.mu.Unlock()
	}
}

func (client *wsClient) broadcast(user *models.User, roomID, content string) error {
	message := Message{
		RoomID:  roomID,
		UserID:  user.ID,
		Content: content,
	}

	roomService := client.handler.services.GetRoomService()
	members, err := roomService.GetRoomMembers(repositories.GetRoomMembersParams{
		RoomID: &roomID,
	}, nil)
	if err != nil {
		return err
	}

	for _, member := range members {
		peer, ok := client.handler.clients[member.UserID]
		if ok {
			err = peer.conn.WriteJSON(message)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
