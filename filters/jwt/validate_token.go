//nolint:revive
package filterjwt

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ===============================================================================================================
// ValidateToken parses and validates a JWT token string using the provided signature key and verification function.
// ===============================================================================================================.
func ValidateToken(tokenType string, jwtSignatureKey []byte) (jwt.MapClaims, string, int, error) {
	troubleShootMsg := ""
	httpCode := http.StatusOK

	// Parse the token string and verify its signing method and signature key
	byteToken, err := jwt.Parse(tokenType, func(jwtToken *jwt.Token) (any, error) {
		if _, successSign := jwtToken.Method.(*jwt.SigningMethodHMAC); !successSign {
			return nil, fmt.Errorf("unexpected signing method %s: %s", tokenType, jwtToken.Header["alg"])
		}

		return jwtSignatureKey, nil
	})
	if err != nil {
		troubleShootMsg = "please ensure the token format is valid or clear the cookies and re-login to get new token"
		httpCode = http.StatusBadRequest

		return nil, troubleShootMsg, httpCode, fmt.Errorf("parse %s token: %w", tokenType, err)
	}

	// Extract and type-assert the token claims as map claims, verifying token validity
	claimToken, successClaim := byteToken.Claims.(jwt.MapClaims)
	if !successClaim || !byteToken.Valid {
		troubleShootMsg = "please ensure the token format is valid or clear the cookies and re-login to get new token"
		httpCode = http.StatusBadRequest

		return nil, troubleShootMsg, httpCode, errors.New("invalid token")
	}

	// Return the validated claim map successfully
	return claimToken, troubleShootMsg, httpCode, nil
}

func CheckTokenExpired(claimToken jwt.MapClaims) (jwt.MapClaims, string, int, error) {
	troubleShootMsg := ""
	httpCode := http.StatusOK

	// Extract and type-assert the expiration timestamp claim from the token claims map
	expiresAt, expiresAtOk := claimToken["expires_at"].(float64)
	if !expiresAtOk {
		troubleShootMsg = "please ensure the token has an expires_at field, or clear the cookies and re-login to get new token"
		httpCode = http.StatusBadRequest

		return nil, troubleShootMsg, httpCode, errors.New("missing expires_at")
	}

	// Verify if the current Unix timestamp has exceeded the token expiration timestamp
	if time.Now().Unix() > int64(expiresAt) {
		troubleShootMsg = "please ensure the token is not expired, or clear the cookies and re-login to get new token"
		httpCode = http.StatusForbidden

		return claimToken, troubleShootMsg, httpCode, fmt.Errorf("token has expired at %f", expiresAt)
	}

	// Return the active unexpired token claims successfully
	return claimToken, troubleShootMsg, httpCode, nil
}
