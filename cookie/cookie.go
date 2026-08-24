package cookie

import (
	"time"
)

// ======================================================================================
// Cookie represents the structure and configuration parameters for an HTTP cookie.
// ======================================================================================.
type Cookie struct {
	// Expires absolute point in time when the cookie expires
	Expires time.Time

	// Name name or key of the cookie
	Name string

	// Value value stored inside the cookie
	Value string

	// Path URL path scope where the cookie is valid
	Path string

	// Domain domain scope where the cookie is valid
	Domain string

	// SameSite cookie SameSite attribute policy (e.g., Lax, Strict, None)
	SameSite string

	// MaxAge maximum age of the cookie in seconds
	MaxAge int

	// Secure restricts cookie transmission to secure (HTTPS) connections only
	Secure bool

	// HTTPOnly prevents client-side scripts from accessing the cookie
	HTTPOnly bool

	// SessionOnly specifies whether the cookie should expire when the browser session ends
	SessionOnly bool
}
