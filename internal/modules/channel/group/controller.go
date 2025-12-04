package group

import (
	"fmt"

	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/logging/log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct {
	group_service *Service
}

func NewController(group_service *Service) *Controller {
	gc := &Controller{group_service: group_service}
	return gc
}

func (gc *Controller) Create(c *fiber.Ctx) error {
	g, err := common.BodyParse[Group](c)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	insert_id, err := gc.group_service.Create(g)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(insert_id)
}

func (gc *Controller) GetAll(c *fiber.Ctx) error {
	gs, err := gc.group_service.GetAll()
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(gs)
}

func (gc *Controller) GetById(c *fiber.Ctx) error {
	id, err := common.ParamsParseUUID(c, "id")
	if err != nil {
		return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
	}

	g, err := gc.group_service.GetByID(id)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(g)
}

func (gc *Controller) GetUsersById(c *fiber.Ctx) error {
	id, err := common.ParamsParseUUID(c, "id")
	if err != nil {
		return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
	}

	us, err := gc.group_service.GetUsers(id)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(us)
}

func (gc *Controller) AddUserToGroup(c *fiber.Ctx) error {
	group_id, err := common.ParamsParseUUID(c, "id")
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
	log.Debug("group_id: %#v", group_id.String())
	log.Debug("body.User_id.String(): %#v", body.User_id.String())

	err = gc.group_service.AddUserToGroup(body.User_id, group_id)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
