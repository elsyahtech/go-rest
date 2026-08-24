//nolint:revive
package filterjwt

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/elsyahtech/go-rest/config"
	"github.com/elsyahtech/go-rest/logger"
)

// ===============================================================================================================
// JWTAuth returns a Fiber middleware handler that validates access and refresh tokens from incoming requests.
// It extracts token claims, populates request locals with token details, and handles unauthorized scenarios.
// ===============================================================================================================.
func JWTAuth() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var (
			err             error
			troubleShootMsg string
			httpCode        int
			tokenJWT        jwt.MapClaims
		)

		// Extract access and refresh tokens from the request context (headers or cookies)
		accessToken, refreshToken := ExtractTokens(ctx)

		// Attempt to validate the access token
		tokenJWT, _, _, err = ValidateAccessToken(ctx, accessToken)
		if err != nil {
			// If access token validation fails, fallback to validating the refresh token
			tokenJWT, troubleShootMsg, httpCode, err = ValidateRefreshToken(ctx, refreshToken)
			if err != nil {
				// Log unsuccessful authentication
				logger.Logger(&logger.LoggerPayload{
					Fields: map[string]any{
						LogFieldKeyTroubleshoot: troubleShootMsg,
						LogFieldKeyHTTPCode:     httpCode,
					},
				}).Error("authentication failed", zap.Error(err))

				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "unauthorized",
				})
			}
		}

		// Extract and type-assert the token issue timestamp claim
		issueAtStr, issueAtOk := tokenJWT["issue_at"].(string)
		if !issueAtOk {
			issueAtStr = time.Now().Format(time.RFC3339)
		}

		// Extract and type-assert the unique token identifier claim
		tokenID, tokenIdOk := tokenJWT["token_id"].(string)
		if !tokenIdOk {
			tokenID = ""
		}

		// Extract and type-assert the token string claim
		token, tokenOk := tokenJWT["token"].(string)
		if !tokenOk {
			token = ""
		}

		// Extract and type-assert the token issuer claim
		issuer, issuerOk := tokenJWT["issuer"].(string)
		if !issuerOk {
			issuer = ""
		}

		// Extract and type-assert the token type claim
		tokenType, tokeTypeOk := tokenJWT["token_type"].(string)
		if !tokeTypeOk {
			tokenType = ""
		}

		// Extract and type-assert the token expiration timestamp claim
		expiresAt, expiresAtOk := tokenJWT["expires_at"].(int64)
		if !expiresAtOk {
			expiresAt = 0
		}

		// Parse the issue time string into a time.Time object, falling back to current time on error
		issueAt, errParse := time.Parse(time.RFC3339, issueAtStr)
		if errParse != nil {
			issueAt = time.Now()
		}

		// Construct the populated token detail structure with extracted claims
		tokenDetail := config.TokenDetail{
			IssueAt:   issueAt,
			TokenID:   tokenID,
			Token:     token,
			Issuer:    issuer,
			Type:      tokenType,
			ExpiresAt: expiresAt,
		}

		// Log successful authentication
		httpCode = http.StatusOK

		logger.Logger(&logger.LoggerPayload{
			Fields: map[string]any{
				LogFieldKeyTokenID:   tokenID,
				LogFieldKeyHTTPCode:  httpCode,
				LogFieldKeyIPAddress: GetClientIP(ctx),
			},
		}).Info("access granted successfully")

		// Store the structured token detail into the Fiber request context locals
		ctx.Locals("tokenJWT", tokenDetail)

		// Proceed to the next middleware or route handler in the chain
		return ctx.Next()
	}
}

func GetClientIP(ctx fiber.Ctx) string {
	// Retrieve the X-Forwarded-For HTTP header which contains proxy forwarding chains
	forwardHeader := ctx.Get("X-Forwarded-For")

	// Extract the first IP address from the comma-separated forwarding list
	firstAddress := strings.Split(forwardHeader, ",")[0]

	// Validate if the trimmed first address is a valid parsed IP address format
	if net.ParseIP(strings.TrimSpace(firstAddress)) != nil {
		return firstAddress
	}

	// Fallback to retrieving the direct connection IP from the Fiber context if forwarding header is invalid or empty
	return ctx.IP()
}
