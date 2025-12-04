package group

import (
	"github.com/NikosGour/chatter/internal/common"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	group_repo Repository
}

func NewController(group_repo Repository) *Controller {
	gc := &Controller{group_repo: group_repo}
	return gc
}

func (gc *Controller) Create(c *fiber.Ctx) error {
	g, err := common.BodyParse[Group](c)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	insert_id, err := gc.group_repo.Create(g)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(insert_id)
}

func (gc *Controller) GetAll(c *fiber.Ctx) error {
	gs, err := gc.group_repo.GetAll()
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

	g, err := gc.group_repo.GetByID(id)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(g)
}
