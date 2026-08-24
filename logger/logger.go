package logger

import (
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/elsyahtech/go-rest/config"
)

var GlobalLogger *zap.Logger

type LoggerPayload struct {
	Fields map[string]any
}

func NewLogger(payload LoggerPayload) *zap.Logger {
	cfg := config.GlobalConfig

	level := zapcore.InfoLevel

	if strings.EqualFold(cfg.App.Mode, "test") {
		level = zapcore.DebugLevel
	}

	configuredLevel := zap.NewAtomicLevelAt(level)

	if !strings.EqualFold(cfg.App.Mode, "production") {
		encoderConfig := newEncoderConfig()
		consoleCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.Lock(os.Stdout),
			configuredLevel,
		)

		return zap.New(consoleCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.Fields(payload.fields()...))
	}

	logDir := filepath.Join(".", "writable", "logs")

	//nolint:mnd
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return zap.NewNop()
	}

	infoFile, closeInfo, err := zap.Open(filepath.Join(logDir, "info.log"))
	if err != nil {
		return zap.NewNop()
	}

	errorFile, _, err := zap.Open(filepath.Join(logDir, "error.log"))
	if err != nil {
		if syncErr := infoFile.Sync(); syncErr != nil {
			_ = syncErr
		}

		closeInfo()

		return zap.NewNop()
	}

	fileEncoder := zapcore.NewJSONEncoder(newEncoderConfig())

	infoLevel := zap.LevelEnablerFunc(func(entryLevel zapcore.Level) bool {
		return entryLevel < zapcore.ErrorLevel && configuredLevel.Enabled(entryLevel)
	})
	errorLevel := zap.LevelEnablerFunc(func(entryLevel zapcore.Level) bool {
		return entryLevel >= zapcore.ErrorLevel && configuredLevel.Enabled(entryLevel)
	})
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, infoFile, infoLevel),
		zapcore.NewCore(fileEncoder, errorFile, errorLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.Fields(payload.fields()...))
}

func (payload LoggerPayload) fields() []zap.Field {
	fields := make([]zap.Field, 0, len(payload.Fields))

	for key, value := range payload.Fields {
		fields = append(fields, zap.Any(key, value))
	}

	return fields
}

func newEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return encoderConfig
}

func Logger(payload *LoggerPayload) *zap.Logger {
	fields := make([]zap.Field, 0, len(payload.Fields))

	fields = append(fields, zap.String("service", config.GlobalConfig.App.AppName))

	for key, value := range payload.Fields {
		fields = append(fields, zap.Any(key, value))
	}

	return GlobalLogger.With(fields...)
}
