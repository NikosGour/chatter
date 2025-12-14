package common

import (
	"fmt"
	"strconv"

	"github.com/NikosGour/logging/log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

type Ctx interface {
	Params(key string, defaultValue ...string) string
}

func ParamsParseInt(c Ctx, field string) (int, error) {
	_v := c.Params(field)
	if _v == "" {
		msg := fmt.Errorf("No %s was provided", field)
		log.Error("%s", msg)
		return 0, msg
	}

	v, err := strconv.Atoi(_v)
	if err != nil {
		msg := fmt.Errorf("Coulnd't convert `%s` to int", _v)
		log.Error("%s", msg)
		return 0, msg
	}

	return v, nil
}

func ParamsParseUUID(c Ctx, field string) (uuid.UUID, error) {
	_v := c.Params(field)
	if _v == "" {
		msg := fmt.Errorf("no %s was provided", field)
		log.Error("%s", msg)
		return uuid.Nil, msg
	}

	id, err := uuid.Parse(_v)
	if err != nil {
		msg := fmt.Errorf("not a valid uuid (%s): `%s`", field, _v)
		log.Error("%s", msg)
		return uuid.Nil, msg
	}

	return id, nil
}
