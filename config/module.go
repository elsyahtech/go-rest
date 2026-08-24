package config

import "github.com/gofiber/fiber/v3"

// ======================================================================================
// ModuleRegistrar defines a function signature for mounting application routes onto Fiber.
// ======================================================================================.
type ModuleRegistrar func(*fiber.App)
