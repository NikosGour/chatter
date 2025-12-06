package user

import (
	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/chatter/internal/modules/channel"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	user_service    *Service
	channel_service *channel.Service
}

func NewController(user_service *Service, channel_service *channel.Service) *Controller {
	uc := &Controller{user_service: user_service, channel_service: channel_service}
	return uc
}

func (uc *Controller) Create(c *fiber.Ctx) error {
	u, err := common.BodyParse[User](c)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	id, err := uc.channel_service.Create(channel.ChannelTypeUser)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}
	u.Id = id

	insert_id, err := uc.user_service.Create(u)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(insert_id)
}

func (uc *Controller) GetAll(c *fiber.Ctx) error {
	us, err := uc.user_service.GetAll()
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(us)
}

func (uc *Controller) GetById(c *fiber.Ctx) error {
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
