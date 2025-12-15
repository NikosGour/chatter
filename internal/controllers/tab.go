package controllers

import (
	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/services"
	"github.com/gofiber/fiber/v2"
)

type TabController struct {
	tab_service *services.TabService
}

func NewTabController(tab_service *services.TabService) *TabController {
	tc := &TabController{tab_service: tab_service}
	return tc
}

func (tc *TabController) Create(c *fiber.Ctx) error {
	tab, err := common.BodyParse[models.Tab](c)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	insert_id, err := tc.tab_service.Create(tab)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(insert_id)
}

func (tc *TabController) GetAll(c *fiber.Ctx) error {
	tabs, err := tc.tab_service.GetAll()
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(tabs)
}

func (tc *TabController) GetById(c *fiber.Ctx) error {
	id, err := common.ParamsParseUUID(c, "id")
	if err != nil {
		return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
	}

	tab, err := tc.tab_service.GetByID(id)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(tab)
}
