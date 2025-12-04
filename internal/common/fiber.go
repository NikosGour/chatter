package common

import (
	"fmt"

	"github.com/NikosGour/logging/log"
	"github.com/gofiber/fiber/v2"
)

func JSONErr(c *fiber.Ctx, message string, status ...int) error {
	if len(status) > 0 {
		return c.Status(status[0]).JSON(fiber.Map{"Error": message})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": message})
}

type Validater interface {
	Validate() error
}

func BodyParse[T Validater](c *fiber.Ctx) (*T, error) {
	v := new(T)
	err := c.BodyParser(v)
	if err != nil {
		msg := fmt.Errorf("on Unmarshal: %w, for body: `%s`", err, c.Body())
		log.Error("%s", msg)
		return nil, msg
	}

	err = (*v).Validate()
	if err != nil {
		msg := fmt.Errorf("on Validate: %w, for body: `%s`", err, c.Body())
		log.Error("%s", msg)
		return nil, msg
	}

	return v, nil
}
