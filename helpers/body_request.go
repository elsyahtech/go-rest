package helpers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

// ======================================================================================
// BodyRequest parses the incoming HTTP request body into the provided struct pointer.
// ======================================================================================.
func (*Helper) BodyRequest(ctx fiber.Ctx, request any) error {
	if err := ctx.Bind().Body(request); err != nil {
		return fmt.Errorf("failed to parse incoming HTTP request body: %w", err)
	}

	return nil
}
