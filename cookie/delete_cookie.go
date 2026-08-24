package cookie

import (
	"github.com/gofiber/fiber/v3"
)

// ======================================================================================
// DeleteCookie removes an HTTP cookie from the client by setting its value to empty and expiring it.
// ======================================================================================.
func DeleteCookie(ctx fiber.Ctx, tokenType string, cookie *Cookie) {
	ctx.Cookie(&fiber.Cookie{
		// Name of the cookie to be deleted
		Name: tokenType,

		// Value is set to empty to clear the cookie contents
		Value: "",

		// Domain scope where the cookie is valid
		Domain: GetDefaultCookiesDomain(cookie.Domain),

		// HTTPOnly flag restricting client-side script access
		HTTPOnly: cookie.HTTPOnly,

		// Secure flag restricting cookie transmission to HTTPS only
		Secure: cookie.Secure,

		// SameSite attribute policy for cross-site request security
		SameSite: GetDefaultCookiesSameSite(cookie.SameSite),

		// Expires timestamp to immediately invalidate the cookie
		Expires: cookie.Expires,
	})
}
