package controllers

import (
	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/services"
	"github.com/gofiber/fiber/v2"
)

type MessageController struct {
	message_service *services.MessageService
}

func NewMessageController(message_service *services.MessageService) *MessageController {
	uc := &MessageController{message_service: message_service}
	return uc
}

func (mc *MessageController) Create(c *fiber.Ctx) error {
	m, err := common.BodyParse[models.Message](c)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	insert_id, err := mc.message_service.Create(m)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(insert_id)
}

func (mc *MessageController) GetAll(c *fiber.Ctx) error {
	messages, err := mc.message_service.GetAll()
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	message_dtos := []services.MessageDTO{}
	for _, message := range messages {
		mdto := mc.message_service.MessageToDTO(&message)
		message_dtos = append(message_dtos, *mdto)
	}

	return c.JSON(message_dtos)
}

func (mc *MessageController) GetById(c *fiber.Ctx) error {
	id, err := common.ParamsParseInt(c, "id")
	if err != nil {
		return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
	}

	message, err := mc.message_service.GetByID(int64(id))
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	mdto := mc.message_service.MessageToDTO(message)

	return c.JSON(mdto)
}
