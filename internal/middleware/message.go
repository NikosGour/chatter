package middleware

// import (
// 	"fmt"

// 	"github.com/NikosGour/chatter/internal/common"
// 	"github.com/NikosGour/logging/log"
// 	"github.com/gofiber/fiber/v2"
// )

// func WithMessangerId(c *fiber.Ctx) error {
// 	cookie := c.Cookies(common.CookieMessangerId)
// 	if cookie == "" {
// 		msg := fmt.Errorf("`%s` cookie is empty", common.CookieMessangerId)
// 		log.Error("%s", msg)
// 		return common.JSONErr(c, msg.Error())
// 	}

// 	// decrypted,err := encryptcookie.DecryptCookie(cookie,common.Dotenv[common.])
// }
