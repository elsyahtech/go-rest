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
// ValidateRefreshToken validates the given refresh token string, verifies its signature,
// checks expiration, and manages expired cookies.
// ======================================================================================.
func ValidateRefreshToken(ctx fiber.Ctx, refreshToken string) (jwt.MapClaims, string, int, error) {
	// Declare variables to hold the parsed token claims and error state
	var (
		err             error
		troubleShootMsg string
		httpCode        int
		validateToken   jwt.MapClaims
	)

	// Validate that the provided access token string is not empty
	if refreshToken == "" {
		troubleShootMsg = "This error occurs in the ValidateRefreshToken method. " +
			"Please ensure the refreshToken object is not empty and check it in ./system/filters/jwt/authentication"
		httpCode = http.StatusBadRequest

		return nil, troubleShootMsg, httpCode, errors.New("missing refresh token")
	}

	cfg := config.GlobalConfig

	const (
		accessTokenName  = "access_token"
		refreshTokenName = "refreh_token"
	)

	// Convert the configured token signature key string into a byte slice
	tokenSignatureKey := []byte(cfg.Token.SignatureKey)

	// Validate the token signature and extract its map claims
	validateToken, troubleShootMsg, httpCode, err = ValidateToken(refreshToken, tokenSignatureKey)
	if err != nil {
		return nil, troubleShootMsg, httpCode, fmt.Errorf("validate refresh token: %w", err)
	}

	// Construct the cookie configuration for the new access token based on global settings
	setUpCookie := cookie.Cookie{
		Expires:     time.Now().Add(time.Duration(cfg.Cookies.MaxAge) * time.Second),
		Name:        accessTokenName,
		Value:       refreshToken,
		Path:        cfg.Cookies.Path,
		Domain:      cfg.Cookies.Domain,
		SameSite:    cfg.Cookies.SameSite,
		MaxAge:      cfg.Cookies.MaxAge,
		Secure:      cfg.Cookies.Secure,
		HTTPOnly:    cfg.Cookies.HTTPOnly,
		SessionOnly: cfg.Cookies.SessionOnly,
	}

	// Check if the token claims indicate that the token has expired
	validateToken, troubleShootMsg, httpCode, err = CheckTokenExpired(validateToken)
	if err != nil {
		setUpCookie.Expires = time.Now().Add(-1 * time.Second)

		cookie.DeleteCookie(ctx, accessTokenName, &setUpCookie)

		return nil, troubleShootMsg, httpCode, fmt.Errorf("refresh token expired: %w, please re-login", err)
	}

	// Extract and type-assert the token ID from the validated claims map
	tokenID, tokenIdOk := validateToken["token_id"].(string)
	if !tokenIdOk {
		tokenID = ""
	}

	// Generate a new access token JWT using the extracted token ID
	_, troubleShootMsg, httpCode, err = GenerateTokenJwt(tokenID, accessTokenName)
	if err != nil {
		return nil, troubleShootMsg, httpCode, fmt.Errorf("generate new access token: %w", err)
	}

	// Create and attach the newly generated access token cookie onto the response context
	troubleShootMsg, httpCode, err = cookie.CreateCookie(ctx, &setUpCookie)
	if err != nil {
		return nil, troubleShootMsg, httpCode, fmt.Errorf("set %s token cookies: %w", accessTokenName, err)
	}

	// Return the validated refresh token map claims successfully along with troubleshooting and HTTP status details
	httpCode = http.StatusOK

	return validateToken, troubleShootMsg, httpCode, nil
}
