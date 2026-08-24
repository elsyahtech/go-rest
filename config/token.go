package config

import (
	"time"
)

// ======================================================================================
// Token represents JWT/OAuth or authentication token generation and lifespan parameters.
// ======================================================================================.
type Token struct {
	// SignatureKey secret key used to sign tokens
	SignatureKey string

	// SigningMethod algorithm method used for signing tokens (e.g., HS256)
	SigningMethod string

	// AccessTokenDuration lifespan duration for access tokens
	AccessTokenDuration time.Duration

	// RefreshTokenDuration lifespan duration for refresh tokens
	RefreshTokenDuration time.Duration
}

// ==============================================================================================================
// TokenDetail represents the detailed metadata and claims associated with a generated authentication token.
// ==============================================================================================================.
type TokenDetail struct {
	IssueAt   time.Time
	TokenID   string
	Token     string
	Issuer    string
	Type      string
	ExpiresAt int64
}

// ======================================================================================
// TokenDefault applies default fallback values to empty fields in the Token configuration.
// ======================================================================================.
func TokenDefault(token Token) *Token {
	// Set default token SignatureKey to 89d142f9593357308c32855027e21dacd92889ff if not specified
	if token.SignatureKey == "" {
		token.SignatureKey = "89d142f9593357308c32855027e21dacd92889ff"
	}

	// Set default token SigningMethod to HS256 if not specified
	if token.SigningMethod == "" {
		token.SigningMethod = "HS256"
	}

	return &token
}

// ==============================================================================================================
// OAuthUserDetail represents the user information extracted from an OAuth2 provider.
// ==============================================================================================================.
type OAuthUserDetail struct {
	RawClaims map[string]any
	IssuedAt  time.Time
	UserID    string
	Email     string
	Name      string
	Provider  string
	ExpiresAt int64
}

// ==============================================================================================================
// AuthPayload represents unified authentication payload for both JWT and OAuth2 flows.
// Stored in ctx.Locals("auth") untuk consistency across auth types.
// ==============================================================================================================.
type AuthPayload struct {
	JWT   *TokenDetail
	OAuth *OAuthUserDetail
	Type  string // "JWT" atau "OAUTH"
}
