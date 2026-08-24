package router

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/elsyahtech/go-rest/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

var GlobalRouter *fiber.App

// ======================================================================================================
// NewRouter creates a new Fiber application instance configured with CORS and registered modules.
// ======================================================================================================.
func NewRouter(routes ...config.ModuleRegistrar) *fiber.App {
	// Initialize a new Fiber instance
	app := fiber.New()

	cfg := config.GlobalConfig.Router

	// Apply CORS middleware using the global router configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		ExposeHeaders:    cfg.ExposeHeaders,
		AllowCredentials: cfg.AllowCredentials,
		AllowOriginsFunc: cfg.AllowOriginsFunc,
		MaxAge:           cfg.MaxAge,
	}))

	// Register each module/route group onto the Fiber application
	for _, register := range routes {
		register(app)
	}

	printStartupBanner(app, routes...)

	return app
}

func printStartupBanner(_ *fiber.App, routes ...config.ModuleRegistrar) {
	cfg := config.GlobalConfig

	address := fmt.Sprintf(
		"%s:%s",
		cfg.Server.Host,
		cfg.Server.Port,
	)

	authType := "JWT"

	if cfg.Filter.AuthenticationType == "OAUTH" {
		authType = "OAUTH (Not support yet)"
	}

	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║ 		             GO REST                             ║")
	fmt.Println("║ 		             v1.0.0                              ║")
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  App Name ...................... %-28s  ║\n", strings.ToUpper(cfg.App.AppName))
	fmt.Printf("║  App Version Running ........... %-28s  ║\n", cfg.App.Version)
	fmt.Printf("║  Host Runnig ................... %-28s  ║\n", address)
	fmt.Printf("║  Database Running .............. %-28s  ║\n", strings.ToUpper(cfg.Database.Driver))
	fmt.Printf("║  Authentication ................ %-28s  ║\n", authType)
	fmt.Printf("║  Modules Running ............... %-28d  ║\n", len(routes))

	for _, register := range routes {
		rawName := moduleName(register)

		nameParts := strings.Split(rawName, ".")
		functionName := nameParts[len(nameParts)-1]
		cleanName := strings.TrimPrefix(functionName, "Register")

		cleanName = strings.TrimSuffix(cleanName, "Routes")

		fmt.Printf("║   - %-58s ║\n", cleanName)
	}

	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func moduleName(register config.ModuleRegistrar) string {
	funFor := runtime.FuncForPC(reflect.ValueOf(register).Pointer())

	if funFor == nil {
		return "unknown"
	}

	return funFor.Name()
}
