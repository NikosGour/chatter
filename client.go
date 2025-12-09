package main

import (
	"github.com/NikosGour/logging/log"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

func main() {
	encryptcookie.GenerateKey()
	log.Debug("encryptcookie.GenerateKey(): %#v", encryptcookie.GenerateKey())
}
