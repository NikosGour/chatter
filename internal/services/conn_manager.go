package services

import (
	"encoding/json"
	"errors"
	"slices"
	"sync"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/logging/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type ConnManager struct {
	clients_mu sync.RWMutex
	Clients    map[uuid.UUID]*websocket.Conn
	broadcast  chan *MessageDTO

	message_service *MessageService
	tab_service     *TabService
	server_service  *ServerService
}

var (
	ErrConnectionNotFound = errors.New("connection not found")
)

func NewConnManager(message_service *MessageService, tab_service *TabService, server_service *ServerService) *ConnManager {
	cm := &ConnManager{
		Clients:         make(map[uuid.UUID]*websocket.Conn),
		broadcast:       make(chan *MessageDTO),
		message_service: message_service,
		tab_service:     tab_service,
		server_service:  server_service,
	}
	return cm
}

func (cm *ConnManager) AddClient(user_id uuid.UUID, conn *websocket.Conn) {
	cm.clients_mu.Lock()
	defer cm.clients_mu.Unlock()

	cm.Clients[user_id] = conn
}

func (cm *ConnManager) RemoveClient(user_id uuid.UUID) error {
	_, err := cm.GetConn(user_id)
	if err != nil {
		return err
	}
	cm.clients_mu.Lock()
	defer cm.clients_mu.Unlock()

	delete(cm.Clients, user_id)

	return nil
}

func (cm *ConnManager) GetConn(user_id uuid.UUID) (*websocket.Conn, error) {
	cm.clients_mu.RLock()
	defer cm.clients_mu.RUnlock()
	conn, ok := cm.Clients[user_id]
	if !ok {
		return nil, ErrConnectionNotFound
	}
	return conn, nil

}

func (cm *ConnManager) ClientReadIncoming(uid uuid.UUID) {
	conn, err := cm.GetConn(uid)
	if err != nil {
		log.Error("on read message: %s", err)
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

			var msg MessageDTO
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
		msg_id, err := cm.message_service.Create(msg)
		if err != nil {
			log.Error("could not insert message to db: %s", err)
			continue
		}
		db_msg, err := cm.message_service.GetByID(msg_id)
		if err != nil {
			log.Error("could not find msg with id: %d, %s", msg_id, err)
		}

		msg_dto := cm.message_service.MessageToDTO(db_msg)
		j_msg, err := json.Marshal(msg_dto)
		if err != nil {
			log.Warn("message failed to be encoded to json: `%#v`, %s", msg, err)
			continue
		}

		tab, err := cm.tab_service.GetByID(msg.Tab.Id)
		if err != nil {
			log.Warn("couldn't find corresponding message tab: `%#v`, %s", msg, err)
		}

		server, err := cm.server_service.GetByID(tab.ServerId)
		if err != nil {
			log.Warn("couldn't find corresponding server: `%#v`, %s", tab, err)
		}

		users, err := cm.server_service.GetUsers(server.Id)
		if err != nil {
			log.Warn("couldn't find users for server: `%#v`, %s", server, err)
		}

		cm.clients_mu.RLock()
		log.Debug("users: %#v", users)
		for uid, conn := range cm.Clients {
			if !slices.ContainsFunc(users, func(u models.User) bool { return u.Id == uid }) {
				continue
			}

			// if !ok {
			// 	log.Warn("tried to write message to offline user: %s", msg.Recipient.GetId())
			// 	continue
			// }

			err := conn.WriteMessage(websocket.TextMessage, j_msg)
			if err != nil {
				log.Error("failed on write: %s", err)
				continue
			}
		}
		cm.clients_mu.RUnlock()
	}
}
