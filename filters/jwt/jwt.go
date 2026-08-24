//nolint:revive
package filterjwt

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elsyahtech/go-rest/config"
	"github.com/golang-jwt/jwt/v5"
)

// ===============================================================================================================
// GenerateTokenJwt generates a signed JSON Web Token (JWT) for either access or refresh token types.
// ===============================================================================================================.
func GenerateTokenJwt(tokenID, tokenType string) (string, string, int, error) {
	troubleShootMsg := ""
	httpCode := http.StatusOK

	const (
		accessToken = "access_token"
		refrehToken = "refresh_token"
	)

	// Declare variables for error handling, token expiration time, and generated token string result
	var (
		err            error
		expiresAt      time.Time
		tokenGenerated string
	)

	// Validate that the token ID parameter is provided and not empty
	if tokenID == "" {
		troubleShootMsg = "please ensure the token object includes a token_id field, or if in login process, please ensure handler success generate token_id before executing the token generation process"
		httpCode = http.StatusBadRequest

		return tokenGenerated, troubleShootMsg, httpCode, errors.New("missing token ID")
	}

	// Validate that the token type parameter is provided and not empty
	if tokenType == "" {
		troubleShootMsg = "please provide a valid token type, ensure token type parameter is set to either 'access_token' or 'refresh_token' before calling Generate Token"
		httpCode = http.StatusBadRequest

		return tokenGenerated, troubleShootMsg, httpCode, errors.New("missing token type")
	}

	// Validate that the token type strictly matches either supported token type constant
	if tokenType != accessToken && tokenType != refrehToken {
		troubleShootMsg = "please ensure the token type passed is strictly 'access_token' or 'refresh_token'"
		httpCode = http.StatusBadRequest

		return tokenGenerated, troubleShootMsg, httpCode, fmt.Errorf("invalid token_type: %s, must be 'access_token' or 'refresh_token'", tokenType)
	}

	cfg := config.GlobalConfig

	// Initialize the token expiration reference starting from the current time
	expiresAt = time.Now()

	// Set the appropriate expiration duration based on whether it is an access or refresh token
	if tokenType == accessToken {
		expiresAt = expiresAt.Add(cfg.Token.AccessTokenDuration)
	} else {
		expiresAt = expiresAt.Add(cfg.Token.RefreshTokenDuration)
	}

	// Create a new JWT token instance utilizing the configured signing method
	tokenSignMethod := jwt.New(GetTokenSigningMethod())

	// Populate the JWT standard and custom claims payload map
	tokenSignMethod.Claims = jwt.MapClaims{
		"issue_at":   time.Now(),
		"token_id":   tokenID,
		"issuer":     cfg.App.AppName,
		"token_type": tokenType,
		"expires_at": expiresAt.Unix(),
	}

	// Sign the token using the secret signature key converted to a byte slice
	tokenGenerated, err = tokenSignMethod.SignedString([]byte(cfg.Token.SignatureKey))
	if err != nil {
		troubleShootMsg = "please check your TokenSignatureKey configuration and TokenSigningMethod in ./app/config/database"
		httpCode = http.StatusInternalServerError

		return tokenGenerated, troubleShootMsg, httpCode, fmt.Errorf("sign token: %w", err)
	}

	// Return successfully with the generated token string
	return tokenGenerated, troubleShootMsg, httpCode, nil
}
