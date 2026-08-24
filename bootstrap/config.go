package bootstrap

import (
	"github.com/elsyahtech/go-rest/config"
)

// ====================================================================
// InitConfig initializes and sets the global system configuration
// ====================================================================.
func InitConfig(provider *config.ConfigProvider) config.Config {
	config.GlobalConfig = config.Config{
		// Load general application configurations
		App: provider.AppProvider(),

		// Load database configurations
		Database: provider.DatabaseProvider(),

		// Load HTTP router configurations
		Router: provider.RouterProvider(),

		// Load Cookies configurations
		Cookies: provider.CookiesProvider(),

		// Load server host, port, etc configurations
		Server: provider.ServerProvider(),

		// Load token configurations
		Token: provider.TokenProvider(),

		// Load middleware and request filter configurations
		Filter: provider.FilterProvider(),

		// Load enabled application modules
		Modules: provider.ModulesProvider(),

		// Load enable MongoDB migrations
		MongoDBMigrations: provider.MongoDBMigrationProvider(),
	}

	return config.GlobalConfig
}
