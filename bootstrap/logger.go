package bootstrap

import (
	"github.com/elsyahtech/go-rest/logger"
	"go.uber.org/zap"
)

// ==================================================================================================
// InitLogger initializes and returns the global logger instance with an empty field map.
// ==================================================================================================.
func InitLogger() *zap.Logger {
	// Create and assign the global logger using the default payload configuration
	logger.GlobalLogger = logger.NewLogger(logger.LoggerPayload{})

	return logger.GlobalLogger
}
