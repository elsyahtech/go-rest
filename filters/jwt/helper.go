//nolint:revive
package filterjwt

import (
	"github.com/elsyahtech/go-rest/config"
	"github.com/golang-jwt/jwt/v5"
)

func GetTokenSigningMethod() jwt.SigningMethod {
	methodName := config.GlobalConfig.Token.SigningMethod
	if methodName == "" {
		return jwt.SigningMethodHS256
	}

	return jwt.GetSigningMethod(methodName)
}
