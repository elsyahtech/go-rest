package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/elsyahtech/go-rest/config"
)

// ==================================================================================================
// RunTesting executes automated unit tests if the application is running in test mode
// ==================================================================================================.
func RunTesting() error {
	appMode := config.GlobalConfig.App.Mode

	if strings.EqualFold(appMode, "test") {
		// Prepare the command to run Go tests without caching and across all subpackages
		command := exec.CommandContext(context.Background(), "go", "test", "-count=1", "./...")

		// Pipe the command's standard output to the console
		command.Stdout = os.Stdout

		// Pipe the command's standard error to the console
		command.Stderr = os.Stderr

		// Execute the test command
		if err := command.Run(); err != nil {
			return fmt.Errorf("error: %w", err)
		}

		return nil
	}

	return nil
}
