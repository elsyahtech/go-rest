package cookie

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

// ======================================================================================
// CreateCookie creates and sets an HTTP cookie on the Fiber response context.
// ======================================================================================.
func CreateCookie(ctx fiber.Ctx, cookie *Cookie) (string, int, error) {
	troubleShootMsg := ""
	httpCode := http.StatusOK

	// Validate that the cookie name is provided
	if cookie.Name == "" {
		troubleShootMsg = "please ensure the cookie name is included to create the cookies"
		httpCode = http.StatusBadRequest

		return troubleShootMsg, httpCode, errors.New("missing cookies name")
	}

	// Validate that the cookie value is provided
	if cookie.Value == "" {
		troubleShootMsg = "please ensure the value is included to create the cookies"
		httpCode = http.StatusBadRequest

		return troubleShootMsg, httpCode, errors.New("token detail is empty")
	}

	// Construct the Fiber cookie configuration using default fallbacks and input parameters
	cookieCreated := &fiber.Cookie{
		Name:        cookie.Name,
		Value:       cookie.Value,
		Domain:      GetDefaultCookiesDomain(cookie.Domain),
		Path:        cookie.Path,
		MaxAge:      cookie.MaxAge,
		Expires:     cookie.Expires,
		Secure:      cookie.Secure,
		HTTPOnly:    cookie.HTTPOnly,
		SameSite:    GetDefaultCookiesSameSite(cookie.SameSite),
		SessionOnly: false,
	}

	// Attach the cookie to the Fiber response context
	ctx.Cookie(cookieCreated)

	return troubleShootMsg, httpCode, nil
}
