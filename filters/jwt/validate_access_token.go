//nolint:revive
package filterjwt

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/go-rest/config"
	"github.com/elsyahtech/go-rest/cookie"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// ======================================================================================
// ValidateAccessToken validates the given access token string, verifies its signature,
// checks expiration, and manages expired cookies.
// ======================================================================================.
func ValidateAccessToken(ctx fiber.Ctx, accessToken string) (jwt.MapClaims, string, int, error) {
	troubleShootMsg := ""

	// Validate that the provided access token string is not empty
	if accessToken == "" {
		troubleShootMsg = "This error occurs in the ValidateAccessToken method. Please ensure the accessToken object is not empty and check it in ./system/filters/jwt/authentication"

		return nil, troubleShootMsg, http.StatusBadRequest, errors.New("missing access token")
	}

	// Define the cookie name constant for access tokens
	const cookieName = "access_token"

	cfg := config.GlobalConfig

	// Convert the configured token signature key string into a byte slice
	tokenSignatureKey := []byte(cfg.Token.SignatureKey)

	// Validate the token signature and extract its map claims
	validateToken, troubleShootMsg, httpCode, err := ValidateToken(accessToken, tokenSignatureKey)
	if err != nil {
		return nil, troubleShootMsg, httpCode, fmt.Errorf("validate access token: %w", err)
	}

	// Check if the token claims indicate that the token has expired
	validateToken, troubleShootMsg, httpCode, err = CheckTokenExpired(validateToken)
	if err != nil {
		cookie.DeleteCookie(ctx, accessToken, &cookie.Cookie{
			Expires:     time.Now().Add(-1 * time.Second),
			Name:        cookieName,
			Value:       accessToken,
			Path:        cfg.Cookies.Path,
			Domain:      cfg.Cookies.Domain,
			SameSite:    cfg.Cookies.SameSite,
			MaxAge:      cfg.Cookies.MaxAge,
			Secure:      cfg.Cookies.Secure,
			HTTPOnly:    cfg.Cookies.HTTPOnly,
			SessionOnly: cfg.Cookies.SessionOnly,
		})
	}

	// Return the validated token map claims
	return validateToken, troubleShootMsg, httpCode, nil
}
