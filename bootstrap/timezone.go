package bootstrap

import (
	"fmt"
	"time"

	"github.com/elsyahtech/go-rest/config"
)

// ==================================================================================================
// InitTimezone loads and sets the local timezone for the application
// ==================================================================================================.
func InitTimezone() error {
	timezone := config.GlobalConfig.App.Timezone

	// Load the location/timezone from the IANA database
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Set the global default timezone for the application
	time.Local = location

	return nil
}
