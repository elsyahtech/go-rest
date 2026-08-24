package config

// ======================================================================================
// GlobalConfig holds the global configuration instance for the entire application.
// ======================================================================================.
var GlobalConfig Config

type ConfigProvider struct {
	AppProvider              func() *App
	DatabaseProvider         func() *Database
	RouterProvider           func() *Router
	CookiesProvider          func() *Cookies
	ServerProvider           func() *Server
	TokenProvider            func() *Token
	FilterProvider           func() *Filter
	ModulesProvider          func() []ModuleRegistrar
	MongoDBMigrationProvider func() []MongoDBMigrationRegistrar
	MongoDBSeederProvider    func() []MongoDBSeederRegistrar
}

// ======================================================================================
// Config represents the master configuration structure containing all subsystem settings.
// ======================================================================================.
type Config struct {
	// App core application settings (mode, timezone, version, name)
	App *App

	// Database connection and migration settings
	Database *Database

	// Router and CORS configuration
	Router *Router

	// Cookies security and scope configuration
	Cookies *Cookies

	// Server network and hosting configuration
	Server *Server

	// Token authentication and lifespan configuration
	Token *Token

	// Filter security or request filtering configuration
	Filter *Filter

	// Modules list of registered application module registrars
	Modules []ModuleRegistrar

	// MongoDBMigrations list of registered MongoDB migrations registrar
	MongoDBMigrations []MongoDBMigrationRegistrar

	// MongoDBSeeders list of registered MongoDB migrations registrar
	MongoDBSeeders []MongoDBSeederRegistrar
}
