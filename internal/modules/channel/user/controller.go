package user

import (
	"github.com/NikosGour/chatter/internal/common"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	user_repo Repository
}

func NewController(user_repo Repository) *Controller {
	uc := &Controller{user_repo: user_repo}
	return uc
}

func (uc *Controller) Create(c *fiber.Ctx) error {
	u, err := common.BodyParse[User](c)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	insert_id, err := uc.user_repo.Create(u)
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(insert_id)
}

func (uc *Controller) GetAll(c *fiber.Ctx) error {
	us, err := uc.user_repo.GetAll()
	if err != nil {
		return common.JSONErr(c, err.Error())
	}

	return c.JSON(us)
}
