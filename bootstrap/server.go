package bootstrap

import (
	"errors"
	"fmt"

	"github.com/elsyahtech/go-rest/config"
	"github.com/elsyahtech/go-rest/router"
	"github.com/gofiber/fiber/v3"
)

// =====================================================================================
// StartServer starts the HTTP server using the global router on the configured port.
// =====================================================================================.
func StartServer() (string, error) {
	troubleShootMsg := ""
	port := config.GlobalConfig.Server.Port

	// Ensure the global router has been properly initialized before starting the server
	if router.GlobalRouter == nil {
		troubleShootMsg = "Ensure router initialization is called before StartServer"

		return troubleShootMsg, errors.New("router is not initialized")
	}

	// Start listening and serving HTTP traffic on the specified port
	if err := router.GlobalRouter.Listen(":"+port, fiber.ListenConfig{
		DisableStartupMessage: true,
	}); err != nil {
		troubleShootMsg = "check server port configuration in ./app/config/server - ensure the port is not already in use"

		return troubleShootMsg, fmt.Errorf("server stopped: %w", err)
	}

	return troubleShootMsg, nil
}
