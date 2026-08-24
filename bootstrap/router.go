package bootstrap

import (
	"errors"

	"github.com/elsyahtech/go-rest/config"
	"github.com/elsyahtech/go-rest/router"
)

// ======================================================================================
// InitRouter initializes and sets up the global HTTP router using registered modules.
// ======================================================================================.
func InitRouter() (string, error) {
	troubleShootMsg := ""
	modules := config.GlobalConfig.Modules

	// Check if any application modules have been registered
	if len(modules) == 0 {
		troubleShootMsg = "ensure all route modules are registered in ./app/config/modules"

		return troubleShootMsg, errors.New("no modules registered")
	}

	// Initialize and assign the global router instance with the registered modules
	router.GlobalRouter = router.NewRouter(modules...)

	return troubleShootMsg, nil
}
