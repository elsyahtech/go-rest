package bootstrap

import (
	"fmt"

	"github.com/elsyahtech/go-rest/database"
)

// ==================================================================================================
// InitDatabase initializes the database connection and assigns it to the global database instance
// ==================================================================================================.
func InitDatabase() (string, error) {
	// Attempt to establish a new database connection and retrieve the action/status context
	dBInitialized, troubleShootMsg, err := database.DBConnection()
	if err != nil {
		return troubleShootMsg, fmt.Errorf("%w", err)
	}

	// Assign the successfully established database connection globally
	database.GlobalDB = dBInitialized

	return "", nil
}
