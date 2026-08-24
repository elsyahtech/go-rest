package config

// ======================================================================================
// Cookies represents default configuration settings for HTTP cookies.
// ======================================================================================.
type Cookies struct {
	// SameSite cookie SameSite attribute policy
	SameSite string

	// Domain cookie domain scope
	Domain string

	// Path cookie path scope
	Path string

	// MaxAge cookie lifetime duration in seconds
	MaxAge int

	// Secure restricts cookie transmission to secure (HTTPS) connections only
	Secure bool

	// HTTPOnly prevents client-side script access to the cookie
	HTTPOnly bool

	// SessionOnly specifies whether the cookie expires when the browser session ends
	SessionOnly bool
}

// ======================================================================================
// CookiesDefault applies default fallback values to empty fields in the Cookies configuration.
// ======================================================================================.
func CookiesDefault(cookies Cookies) *Cookies {
	// Set default cookies samesite to development if not specified
	if cookies.SameSite == "" {
		cookies.SameSite = "development"
	}

	// Set default cookies domain to localhost if not specified
	if cookies.Domain == "" {
		cookies.Domain = "localhost"
	}

	// Set default cookies path to / if not specified
	if cookies.Path == "" {
		cookies.Path = "/"
	}

	// Set default cookies maxAge to 63072000 if not specified
	if cookies.MaxAge == 0 {
		cookies.MaxAge = 63072000
	}

	return &cookies
}
