package internal

import (
	"fmt"

	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/chatter/internal/controllers"
	"github.com/NikosGour/chatter/internal/repositories"
	"github.com/NikosGour/chatter/internal/services"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/NikosGour/logging/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

type APIServer struct {
	listening_addr string
	db             *storage.PostgreSQLStorage

	user_controller    *controllers.UserController
	server_controller  *controllers.ServerController
	message_controller *controllers.MessageController
	tab_controller     *controllers.TabController

	conn_manager *services.ConnManager
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

	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | Params: ${queryParams} | ReqBody: ${body} | ResBody: ${resBody} | ${error}\n",
	}))

	s.DependencyInjection()

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

		err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "nikos"))
		if err != nil {
			log.Error("on WriteMessage: %s", err)
			return
		}
		log.Info("Sent")
	}))

	ws.Get("/messages", websocket.New(func(c *websocket.Conn) {

		defer func() {
			err := c.Close()
			if err != nil {
				log.Warn("on conn close path: /ws/messages got error: %s", err)
			}
		}()

		_id := c.Query("uid")
		id, err := uuid.Parse(_id)
		if err != nil {
			log.Error("%s", err)
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseUnsupportedData, fmt.Sprintf("uid param is not a valid uuid: `%s`", _id)))
			if err != nil {
				log.Error("couldn't inform the client of malformed uuid: %s", err)
			}
			return
		}
		s.conn_manager.AddClient(id, c)
		log.Debug("%#v", s.conn_manager.Clients)
		defer s.conn_manager.RemoveClient(id)

		s.conn_manager.ClientReadIncoming(id)
	}))

	user := app.Group("/user")
	user.Post("/", s.user_controller.Create)
	user.Get("/", s.user_controller.GetAll)
	user.Get("/:id", s.user_controller.GetById)

	group := app.Group("/server")
	group.Post("/", s.server_controller.Create)
	group.Get("/", s.server_controller.GetAll)
	group.Get("/:id", s.server_controller.GetById)
	group.Get("/:id/users", s.server_controller.GetUsersById)
	group.Get("/:id/tabs", s.server_controller.GetTabsById)
	group.Post("/:id", s.server_controller.AddUserToServer)

	message := app.Group("/message")
	message.Post("/", s.message_controller.Create)
	message.Get("/", s.message_controller.GetAll)
	message.Get("/:id", s.message_controller.GetById)

	tab := app.Group("/tab")
	tab.Post("/", s.tab_controller.Create)
	tab.Get("/", s.tab_controller.GetAll)
	tab.Get("/:id", s.tab_controller.GetById)

	return app
}

func (s *APIServer) DependencyInjection() {

	user_repo := repositories.NewUserRepository(s.db)
	server_repo := repositories.NewServerRepository(s.db)
	message_repo := repositories.NewMessageRepository(s.db)
	tab_repo := repositories.NewTabRepository(s.db)

	user_service := services.NewUserService(user_repo)
	message_service := services.NewMessageService(message_repo)
	tab_service := services.NewTabService(tab_repo)
	server_service := services.NewServerService(server_repo, user_service, tab_service)

	s.conn_manager = services.NewConnManager(message_service)
	go s.conn_manager.HandleIncomingMessages()

	s.user_controller = controllers.NewUserController(user_service)
	s.server_controller = controllers.NewServerController(server_service)
	s.message_controller = controllers.NewMessageController(message_service)
	s.tab_controller = controllers.NewTabController(tab_service)
}
