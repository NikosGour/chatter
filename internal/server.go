package internal

import (
	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/chatter/internal/modules/channel"
	"github.com/NikosGour/chatter/internal/modules/channel/group"
	"github.com/NikosGour/chatter/internal/modules/channel/user"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/NikosGour/logging/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type APIServer struct {
	listening_addr string
	db             *storage.PostgreSQLStorage

	user_controller  *user.Controller
	group_controller *group.Controller
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
		log.Info("Sent")
	}))

	ws.Get("/message", func(c *fiber.Ctx) error { return nil })

	user := app.Group("/user")
	user.Post("/", s.user_controller.Create)
	user.Get("/", s.user_controller.GetAll)
	user.Get("/:id", s.user_controller.GetById)

	group := app.Group("/group")
	group.Post("/", s.group_controller.Create)
	group.Get("/", s.group_controller.GetAll)
	group.Get("/:id", s.group_controller.GetById)
	group.Post("/:id", s.group_controller.AddUserToGroup)
	return app
}

func (s *APIServer) DependencyInjection() {
	channel_repo := channel.NewRepository(s.db)
	user_repo := user.NewRepository(s.db, channel_repo)
	group_repo := group.NewRepository(s.db, channel_repo, user_repo)
	s.user_controller = user.NewController(user_repo)
	s.group_controller = group.NewController(group_repo)
}
