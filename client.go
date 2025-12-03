package main

import (
	"io"

	"github.com/NikosGour/logging/log"
	"github.com/fasthttp/websocket"
)

func main() {
	d := new(websocket.Dialer)
	conn, res, err := d.Dial("ws://localhost:8080", nil)
	if err != nil {
		log.Fatal("on Dial: %s", err)
	}
	defer conn.Close()
	log.Debug("res: %#v", res)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("on ReadAll: %s", err)
	}
	log.Debug("b: %#v", b)

	mt, m, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("on ReadMessage: %s", err)
	}
	log.Debug("mt: %#v", mt)
	log.Debug("m: %#v", string(m))
}
