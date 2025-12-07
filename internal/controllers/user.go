package controllers

import (
	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/services"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	user_service *services.UserService
}

func NewUserController(user_service *services.UserService) *UserController {
	uc := &UserController{user_service: user_service}
	return uc
}

func (uc *UserController) Create(c *fiber.Ctx) error {
	u, err := common.BodyParse[models.User](c)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	insert_id, err := uc.user_service.Create(u)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(insert_id)
}

func (uc *UserController) GetAll(c *fiber.Ctx) error {
	us, err := uc.user_service.GetAll()
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(us)
}

func (uc *UserController) GetById(c *fiber.Ctx) error {
	id, err := common.ParamsParseUUID(c, "id")
	if err != nil {
		return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
	}

	u, err := uc.user_service.GetByID(id)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(u)
}
