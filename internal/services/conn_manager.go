package services

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/logging/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type ConnManager struct {
	clients_mu sync.RWMutex
	clients    map[uuid.UUID]*websocket.Conn
	broadcast  chan *models.MessageDTO

	message_service *MessageService
}

var (
	ErrConnectionNotFound = errors.New("connection not found")
)

func NewConnManager(message_service *MessageService) *ConnManager {
	cm := &ConnManager{
		clients:         make(map[uuid.UUID]*websocket.Conn),
		broadcast:       make(chan *models.MessageDTO),
		message_service: message_service,
	}
	return cm
}

func (cm *ConnManager) AddClient(id uuid.UUID, conn *websocket.Conn) {
	cm.clients_mu.Lock()
	defer cm.clients_mu.Unlock()

	cm.clients[id] = conn
}

func (cm *ConnManager) RemoveClient(id uuid.UUID) error {
	_, err := cm.GetConn(id)
	if err != nil {
		return err
	}
	cm.clients_mu.Lock()
	defer cm.clients_mu.Unlock()

	delete(cm.clients, id)

	return nil
}

func (cm *ConnManager) GetConn(id uuid.UUID) (*websocket.Conn, error) {
	cm.clients_mu.RLock()
	conn, ok := cm.clients[id]
	cm.clients_mu.RUnlock()

	if !ok || conn == nil {
		return nil, ErrConnectionNotFound
	}
	return conn, nil
}

func (cm *ConnManager) ClientReadIncoming(id uuid.UUID) {
	conn, err := cm.GetConn(id)
	if err != nil {
		log.Error("getConn: %s", err)
		return
	}

	for {
		mt, data, err := conn.ReadMessage()
		if err != nil {
			// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Error("on read message: %s", err)
			return
			// }
		}
		if mt > 0 {
			log.Debug("mt: %#v ,data: %#v", mt, string(data))

			var msg models.MessageDTO
			err := json.Unmarshal(data, &msg)
			if err != nil {
				log.Error("on unmarshal: %s", err)
				continue
			}

			cm.broadcast <- &msg
		}
	}
}

func (cm *ConnManager) HandleIncomingMessages() {
	for msg := range cm.broadcast {
		// cm.clients_mu.RLock()
		// conn, ok := cm.clients[msg.Recipient.GetId()]
		// cm.clients_mu.RUnlock()

		for _, conn := range cm.clients {

			// msg_id, err := cm.message_service.Create(msg)
			// if err != nil {
			// 	log.Error("could not insert message to db: %s", err)
			// 	continue
			// }
			// msg.Id = msg_id

			// if !ok {
			// 	log.Warn("tried to write message to offline user: %s", msg.Recipient.GetId())
			// 	continue
			// }

			err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Text))
			if err != nil {
				log.Error("failed on write: %s", err)
				continue
			}
		}
	}
}
