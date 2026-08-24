//nolint:revive
package filterjwt

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

// ===============================================================================================================
// ExtractTokens extracts the access and refresh tokens from the incoming HTTP request headers or cookies.
// ===============================================================================================================.
func ExtractTokens(ctx fiber.Ctx) (string, string) {
	// Retrieve the Authorization header value from the request context
	authorization := ctx.Get("Authorization")

	// Initialize variables to hold the extracted access and refresh tokens
	var accessToken, refreshToken string

	// Check if the Authorization header contains a Bearer token scheme
	if strings.HasPrefix(authorization, "Bearer ") {
		// Trim the "Bearer " prefix to isolate the raw token string
		token := strings.TrimPrefix(authorization, "Bearer ")

		// Assign the token to the access token variable
		accessToken = token

		// Assign the token to the refresh token variable as well for fallback or dual usage
		refreshToken = token
	} else {
		// Fallback to extracting the access token from request cookies
		accessToken = ctx.Cookies("access_token")

		// Fallback to extracting the refresh token from request cookies
		refreshToken = ctx.Cookies("refresh_token")
	}

	// Return the extracted access and refresh tokens
	return accessToken, refreshToken
}
