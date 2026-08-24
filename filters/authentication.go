package filters

import (
	"github.com/elsyahtech/go-rest/config"
	filterjwt "github.com/elsyahtech/go-rest/filters/jwt"
	"github.com/gofiber/fiber/v3"
)

func Authentication() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authType := config.GlobalConfig.Filter.AuthenticationType

		//nolint:revive,nolintlint
		switch authType {
		case "JWT":
			return filterjwt.JWTAuth()(ctx)
		case "OAUTH":
			return filterjwt.JWTAuth()(ctx)
		default:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unknown authentication type",
			})
		}
	}
}
