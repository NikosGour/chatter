package channel

// import (
// 	"fmt"

// 	"github.com/NikosGour/chatter/internal/common"
// 	"github.com/NikosGour/chatter/internal/modules/channel/group"
// 	"github.com/NikosGour/chatter/internal/modules/channel/user"
// 	"github.com/NikosGour/logging/log"
// 	"github.com/gofiber/fiber/v2"
// )

// type Controller struct {
// 	channel_service *Service
// }

// func NewController(channel_service *Service) *Controller {
// 	uc := &Controller{channel_service: channel_service}
// 	return uc
// }

// func (uc *Controller) Create(c *fiber.Ctx) error {
// 	chtype, err := common.ParamsParseString(c, "chtype")
// 	if err != nil {
// 		return common.JSONErr(c, err.Error())
// 	}

// 	var ch Channel
// 	switch ChannelType(chtype) {
// 	case ChannelTypeUser:
// 		user, err := common.BodyParse[user.User](c)
// 		if err != nil {
// 			return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
// 		}

// 		ch = user
// 	case ChannelTypeGroup:
// 		group, err := common.BodyParse[group.Group](c)
// 		if err != nil {
// 			return common.JSONErr(c, err.Error(), fiber.StatusBadRequest)
// 		}

// 		ch = group
// 	default:
// 		msg := fmt.Errorf("invalid ChannelType provided: `%s`", chtype)
// 		log.Error("%s", msg)
// 		return common.JSONErr(c, msg.Error(), fiber.StatusBadRequest)
// 	}

// 	id, err := uc.channel_service.Create(ChannelType(chtype), ch)
// 	if err != nil {
// 		return common.JSONErr(c, err.Error())
// 	}

// 	return c.JSON(id)
// }
