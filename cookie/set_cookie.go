package cookie

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/go-rest/config"
	"github.com/gofiber/fiber/v3"
)

// ====================================================================================================================
// SetCookie creates and sets a single HTTP cookie with specified name and token value using global configurations.
// ====================================================================================================================.
func SetCookie(ctx fiber.Ctx, cookieName, token string) (string, int, error) {
	troubleShootMsg := ""

	var err error

	// Validate that the cookie name is provided
	if cookieName == "" {
		troubleShootMsg = "please ensure the cookie name is included to set the cookie"

		return troubleShootMsg, http.StatusOK, errors.New("missing token type")
	}

	// Validate that the token value is provided
	if token == "" {
		troubleShootMsg = "please ensure the token is included to set the cookie"

		return troubleShootMsg, http.StatusOK, errors.New("missing token type")
	}

	cfg := config.GlobalConfig

	// Build the cookie configuration using global settings and computed expiration time
	cookie := Cookie{
		Expires:     time.Now().Add(time.Duration(cfg.Cookies.MaxAge) * time.Second),
		Name:        cookieName,
		Value:       token,
		Path:        cfg.Cookies.Path,
		Domain:      cfg.Cookies.Domain,
		SameSite:    cfg.Cookies.SameSite,
		MaxAge:      cfg.Cookies.MaxAge,
		Secure:      cfg.Cookies.Secure,
		HTTPOnly:    cfg.Cookies.HTTPOnly,
		SessionOnly: cfg.Cookies.SessionOnly,
	}

	// Invoke CreateCookie to attach the single cookie onto the response context
	troubleShootMsg, httpCode, err := CreateCookie(ctx, &cookie)
	if err != nil {
		return troubleShootMsg, httpCode, fmt.Errorf("set cookie: %w", err)
	}

	// Return
	return troubleShootMsg, http.StatusOK, nil
}

// =============================================================================================================
// SetCookies sets multiple authentication cookies (access token and refresh token) on the response context
// =============================================================================================================.
func SetCookies(ctx fiber.Ctx, accessToken, refreshToken *Cookie) (string, int, error) {
	troubleShootMsg := ""
	httpCode := http.StatusOK

	var (
		err    error
		cookie Cookie
	)

	// Define a slice of anonymous structures mapping cookie objects to their respective names
	cookies := []struct {
		Value *Cookie
		Name  string
	}{
		{
			accessToken,
			"access_token",
		},
		{
			refreshToken,
			"refresh_token",
		},
	}

	// Iterate through each cookie definition and create them sequentially on the context
	for _, c := range cookies {
		cookie.Name = c.Name
		cookie.Value = c.Value.Value

		troubleShootMsg, httpCode, err = CreateCookie(ctx, &cookie)
		if err != nil {
			return troubleShootMsg, httpCode, fmt.Errorf("%w", err)
		}
	}

	// Return
	return troubleShootMsg, httpCode, nil
}
