package internal

import (
	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/NikosGour/logging/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type APIServer struct {
	listening_addr string
	db             *storage.PostgreSQLStorage
}

func NewAPIServer(db *storage.PostgreSQLStorage) *APIServer {
	s := &APIServer{db: db}
	listening_addr := common.Dotenv[common.EnvHOST_ADDRESS] + ":" + common.Dotenv[common.EnvPORT]
	s.listening_addr = listening_addr
	return s
}

func (s *APIServer) Start() {
	app := s.SetupServer()

	err := app.Listen(s.listening_addr)
	if err != nil {
		log.Fatal("%s", err)
	}
}

func (s *APIServer) SetupServer() *fiber.App {
	app := fiber.New()

	ws := app.Group("/ws", func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}

		return c.Next()
	})

	ws.Get("/test", websocket.New(func(c *websocket.Conn) {
		err := c.WriteMessage(websocket.TextMessage, []byte("nikos"))
		if err != nil {
			log.Error("on WriteMessage: %s", err)
			return
		}
		log.Info("Sent")
	}))

	ws.Get("/message", func(c *fiber.Ctx) error { return nil })
	return app
}
