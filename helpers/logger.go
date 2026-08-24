package helpers

import (
	"github.com/elsyahtech/go-rest/logger"
)

func (*Helper) LoggerError(message string, fields map[string]any) {
	logFields := fields

	if logFields == nil {
		logFields = make(map[string]any)
	}

	logger.Logger(&logger.LoggerPayload{
		Fields: logFields,
	}).Error(message)
}

func (*Helper) LoggerInfo(message string, fields map[string]any) {
	logFields := fields

	if logFields == nil {
		logFields = make(map[string]any)
	}

	logger.Logger(&logger.LoggerPayload{
		Fields: logFields,
	}).Info(message)
}
