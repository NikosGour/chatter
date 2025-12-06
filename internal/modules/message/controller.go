package message

import (
	"github.com/NikosGour/chatter/internal/common"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	message_service *Service
}

func NewController(message_service *Service) *Controller {
	uc := &Controller{message_service: message_service}
	return uc
}

func (mc *Controller) Create(c *fiber.Ctx) error {
	m, err := common.BodyParse[Message](c)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	insert_id, err := mc.message_service.Create(m)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(insert_id)
}

func (mc *Controller) GetAll(c *fiber.Ctx) error {
	messages, err := mc.message_service.GetAll()
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(messages)
}

func (mc *Controller) GetById(c *fiber.Ctx) error {
	id, err := common.ParamsParseInt(c, "id")
	if err != nil {
		return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
	}

	message, err := mc.message_service.GetByID(int64(id))
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(message)
}
