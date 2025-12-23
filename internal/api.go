package internal

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/chatter/internal/controllers"
	"github.com/NikosGour/chatter/internal/models"
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

	user_service    *services.UserService
	server_service  *services.ServerService
	message_service *services.MessageService
	tab_service     *services.TabService

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
	s.SetupDummyData()

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
	message.Get("/tab/:tab_id", s.message_controller.GetByTabId)

	tab := app.Group("/tab")
	tab.Post("/", s.tab_controller.Create)
	tab.Get("/", s.tab_controller.GetAll)
	tab.Get("/:id", s.tab_controller.GetById)

	return app
}

func (s *APIServer) DependencyInjection() {

	user_repo := repositories.NewUserRepository(s.db)
	tab_repo := repositories.NewTabRepository(s.db)
	message_repo := repositories.NewMessageRepository(s.db)
	server_repo := repositories.NewServerRepository(s.db)

	s.user_service = services.NewUserService(user_repo)
	s.tab_service = services.NewTabService(tab_repo)
	s.message_service = services.NewMessageService(message_repo, s.tab_service)
	s.server_service = services.NewServerService(server_repo, s.user_service, s.tab_service)

	s.conn_manager = services.NewConnManager(s.message_service, s.tab_service, s.server_service)
	go s.conn_manager.HandleIncomingMessages()

	s.user_controller = controllers.NewUserController(s.user_service)
	s.tab_controller = controllers.NewTabController(s.tab_service)
	s.message_controller = controllers.NewMessageController(s.message_service)
	s.server_controller = controllers.NewServerController(s.server_service)
}

func (s *APIServer) SetupDummyData() {
	usernames := []string{"nikos", "maria", "rinos", "nisfa", "gkai", "mitsos"}
	user_ids := []uuid.UUID{}
	for _, username := range usernames {

		user, err := s.user_service.GetByTestUsername(username)
		if err != nil {
			if !errors.Is(err, models.ErrUserNotFound) {
				log.Warn("%s", err)
				continue
			}

		}
		if len(user) != 0 {
			user_ids = append(user_ids, user[0].Id)
			continue
		}

		id, err := s.user_service.Create(&repositories.UserDBO{Username: username, Password: "123", DateCreated: time.Now(), IsTest: true})
		if err != nil {
			log.Warn("%s", err)
		}
		user_ids = append(user_ids, id)
	}

	servers := []struct {
		Name     string
		User_ids []int
	}{
		{
			Name:     "Gamiades",
			User_ids: []int{0, 1, 2},
		},
		{
			Name:     "HUA",
			User_ids: []int{3, 4, 5},
		},
		{
			Name:     "CTF",
			User_ids: []int{1, 3, 5},
		},
		{
			Name:     "ArchUsers",
			User_ids: []int{0, 2, 4},
		},
	}
	server_ids := []uuid.UUID{}
	for _, ser := range servers {
		server, err := s.server_service.GetByTestName(ser.Name)
		if err != nil {
			if !errors.Is(err, models.ErrServerNotFound) {
				log.Warn("%s", err)
				continue
			}
		}
		if len(server) != 0 {
			server_ids = append(server_ids, server[0].Id)
			continue
		}

		id, err := s.server_service.Create(&repositories.ServerDBO{Name: ser.Name, DateCreated: time.Now(), IsTest: true})
		if err != nil {
			log.Warn("%s", err)
		}
		server_ids = append(server_ids, id)

		for _, idx := range ser.User_ids {
			user_id := user_ids[idx]
			err := s.server_service.AddUserToServer(user_id, id)
			if err != nil {
				log.Warn("%s", err)
			}
		}
	}

	tabs := []struct {
		Name        string
		Servers_idx []int
	}{
		{Name: "General", Servers_idx: []int{0, 1, 2}},
		{Name: "Memes", Servers_idx: []int{0, 1, 2, 3}},
		{Name: "Studying", Servers_idx: []int{2, 3}},
	}

	for _, t := range tabs {
		_, err := s.tab_service.GetByName(t.Name)
		if err != nil {
			if errors.Is(err, models.ErrTabNotFound) {
				continue
			}
			log.Warn("%s", err)
			continue

		}

		for _, idx := range t.Servers_idx {
			server_id := server_ids[idx]
			tabs, err := s.server_service.GetTabs(server_id)
			if err != nil {
				log.Warn("%s", err)
				continue
			}

			if slices.ContainsFunc(tabs, func(tab models.Tab) bool {
				return tab.Name == t.Name
			}) {
				continue
			}

			_, err = s.tab_service.Create(&repositories.TabDBO{Name: t.Name, ServerId: server_id, DateCreated: time.Now()})
			if err != nil {
				log.Warn("%s", err)
			}

		}
	}
}
