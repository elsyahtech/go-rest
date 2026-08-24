package helpers

import (
	"fmt"

	"github.com/elsyahtech/go-rest/view"
	"github.com/gofiber/fiber/v3"
)

func (*Helper) JSON(fields map[string]any) *view.Response {
	viewFields := fields

	if viewFields == nil {
		viewFields = make(map[string]any)
	}

	return view.View(&view.JSON{
		Fields: viewFields,
	})
}

//nolint:staticcheck,revive
func (*Helper) View(ctx fiber.Ctx, response *view.Response, err error) error {
	if err = ctx.Status(response.HTTPCode).JSON(response); err != nil {
		return fmt.Errorf("failed to send JSON Response: %w", err)
	}

	return nil
}
