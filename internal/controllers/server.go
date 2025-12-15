package controllers

import (
	"fmt"

	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/services"
	"github.com/NikosGour/logging/log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ServerController struct {
	server_service *services.ServerService
}

func NewServerController(server_service *services.ServerService) *ServerController {
	sc := &ServerController{server_service: server_service}
	return sc
}

func (sc *ServerController) Create(c *fiber.Ctx) error {
	server, err := common.BodyParse[models.Server](c)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	insert_id, err := sc.server_service.Create(server)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(insert_id)
}

func (sc *ServerController) GetAll(c *fiber.Ctx) error {
	servers, err := sc.server_service.GetAll()
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(servers)
}

func (sc *ServerController) GetById(c *fiber.Ctx) error {
	id, err := common.ParamsParseUUID(c, "id")
	if err != nil {
		return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
	}

	g, err := sc.server_service.GetByID(id)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(g)
}

func (sc *ServerController) GetUsersById(c *fiber.Ctx) error {
	id, err := common.ParamsParseUUID(c, "id")
	if err != nil {
		return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
	}

	users, err := sc.server_service.GetUsers(id)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(users)
}

func (sc *ServerController) GetTabsById(c *fiber.Ctx) error {
	id, err := common.ParamsParseUUID(c, "id")
	if err != nil {
		return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
	}

	tabs, err := sc.server_service.GetTabs(id)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(tabs)
}

func (sc *ServerController) AddUserToServer(c *fiber.Ctx) error {
	server_id, err := common.ParamsParseUUID(c, "id")
	if err != nil {
		return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
	}

	body := &struct {
		User_id uuid.UUID `json:"user_id"`
	}{}
	err = c.BodyParser(body)
	if err != nil {
		msg := fmt.Errorf("on Unmarshal: %w, for body: `%s`", err, c.Body())
		log.Error("%s", msg)
		return common.JSONErr(c, msg.Error(), fiber.StatusBadRequest)
	}
	log.Debug("server_id: %#v", server_id.String())
	log.Debug("body.User_id.String(): %#v", body.User_id.String())

	err = sc.server_service.AddUserToServer(body.User_id, server_id)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
