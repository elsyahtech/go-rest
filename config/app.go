package config

// ======================================================================================
// App represents general application-level metadata and runtime options.
// ======================================================================================.
type App struct {
	// Application runtime mode (e.g., development, production, test)
	Mode string

	// Application default timezone
	Timezone string

	// Application version string
	Version string

	// AppName display name
	AppName string
}

// ======================================================================================
// AppDefault applies default fallback values to empty fields in the App configuration.
// ======================================================================================.
func AppDefault(app App) *App {
	// Set default application mode to development if not specified
	if app.Mode == "" {
		app.Mode = "development"
	}

	// Set default timezone to UTC if not specified
	if app.Timezone == "" {
		app.Timezone = "UTC"
	}

	// Set default version to 1.0.0 if not specified
	if app.Version == "" {
		app.Version = "1.0.0"
	}

	// Set default application name if not specified
	if app.AppName == "" {
		app.AppName = "Golang Rest Api Starter"
	}

	return &app
}
