package token

import (
	"fmt"

	"github.com/elsyahtech/go-rest/config"
	filterjwt "github.com/elsyahtech/go-rest/filters/jwt"
)

func GenerateToken(tokenID, tokenType string) (token string, troubleShootMsg string, httpCode int, err error) {
	authenticationType := config.GlobalConfig.Filter.AuthenticationType

	//nolint:revive,nolintlint
	switch authenticationType {
	case "JWT":
		token, troubleShootMsg, httpCode, err = filterjwt.GenerateTokenJwt(tokenID, tokenType)
		if err != nil {
			return token, troubleShootMsg, httpCode, fmt.Errorf("%w", err)
		}
	default:
		token, troubleShootMsg, httpCode, err = filterjwt.GenerateTokenJwt(tokenID, tokenType)
		if err != nil {
			return token, troubleShootMsg, httpCode, fmt.Errorf("%w", err)
		}
	}

	return token, troubleShootMsg, httpCode, nil
}
