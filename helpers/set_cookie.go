package helpers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/elsyahtech/go-rest/cookie"
	"github.com/gofiber/fiber/v3"
)

type Cookie = cookie.Cookie

func (*Helper) SetCookie(ctx fiber.Ctx, cookieName, token string) (troubleShootMsg string, httpCode int, err error) {
	if cookieName == "" || token == "" {
		return "ensure cookie name cannot be empty", http.StatusBadRequest, errors.New("cookie name is empty")
	}

	troubleShootMsg, httpCode, err = cookie.SetCookie(ctx, cookieName, token)
	if err != nil {
		return troubleShootMsg, httpCode, fmt.Errorf("%w", err)
	}

	return troubleShootMsg, httpCode, nil
}
