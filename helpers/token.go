package helpers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"

	"github.com/elsyahtech/go-rest/config"
	"github.com/elsyahtech/go-rest/token"
)

func (*Helper) GetToken(ctx fiber.Ctx) (result *config.TokenDetail, err error) {
	authenticationType := config.GlobalConfig.Filter.AuthenticationType

	var (
		tokenDetail config.TokenDetail
		isSuccess   bool
	)

	//nolint:revive
	switch authenticationType {
	case "JWT":
		tokenDetail, isSuccess = ctx.Locals("tokenJWT").(config.TokenDetail)
		if !isSuccess {
			return nil, errors.New("user unauthorized")
		}
	case "OAUTH":
		tokenDetail, isSuccess = ctx.Locals("tokenJWT").(config.TokenDetail)
		if !isSuccess {
			return nil, errors.New("user unauthorized")
		}
	default:
		tokenDetail, isSuccess = ctx.Locals("tokenJWT").(config.TokenDetail)
		if !isSuccess {
			return nil, errors.New("user unauthorized")
		}
	}

	return &tokenDetail, nil
}

func (*Helper) GenerateTokenAuthentication(tokenID, tokenType string) (
	tokenGenerated string, troubleShootMsg string, httpCode int, err error,
) {
	tokenGenerated, troubleShootMsg, httpCode, err = token.GenerateToken(tokenID, tokenType)
	if err != nil {
		return "", troubleShootMsg, httpCode, fmt.Errorf("%w", err)
	}

	return tokenGenerated, troubleShootMsg, httpCode, nil
}
