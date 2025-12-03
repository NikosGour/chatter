package main

import (
	log "github.com/NikosGour/logging/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", websocket.New(func(c *websocket.Conn) {
		err := c.WriteMessage(websocket.TextMessage, []byte("nikos"))
		if err != nil {
			log.Error("on WriteMessage: %s", err)
			return
		}

		log.Info("Sent")
	}))

	log.Fatal("%s", app.Listen(":8080"))
}
