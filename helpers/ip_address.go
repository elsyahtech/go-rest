package helpers

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// ======================================================================================
// GetClientIP extracts and returns the real client IP address from the request headers
// or fallback context.
// ======================================================================================.
func (*Helper) GetClientIP(ctx fiber.Ctx) string {
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
