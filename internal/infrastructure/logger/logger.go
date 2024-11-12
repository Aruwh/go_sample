package logger

import (
	"context"
	"fewoserv/internal/infrastructure/common"
	"fewoserv/internal/infrastructure/config"
	"fewoserv/internal/infrastructure/utils"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger struct {
		logLevel common.LogLevel
		label    string
		_logger  *zap.Logger
	}
)

var configZap = zap.Config{
	OutputPaths: []string{"stdout"},
	Encoding:    "json",
	EncoderConfig: zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	},
}

func getZapLogLvl(level common.LogLevel) zapcore.Level {
	var zapLogLvl = zapcore.DebugLevel

	switch level {
	case common.DEBUG:
		zapLogLvl = zapcore.DebugLevel
	case common.INFO:
		zapLogLvl = zapcore.InfoLevel
	case common.WARN:
		zapLogLvl = zapcore.WarnLevel
	case common.ERROR:
		zapLogLvl = zapcore.ErrorLevel
	case common.PANIC:
		zapLogLvl = zapcore.PanicLevel
	}

	return zapLogLvl
}

// New creates a new Logger instance with the specified log level.
func New(label string) *Logger {
	cfg := config.Load()

	// Öffnen Sie die Protokolldatei im Append-Modus. Sie können dies anpassen.
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	configZap.Level = zap.NewAtomicLevelAt(getZapLogLvl(cfg.Service.LogLevel))
	configZap.OutputPaths = []string{"stdout", "log.txt"}
	configZap.ErrorOutputPaths = []string{"stdout", "log.txt"}

	_logger, err := configZap.Build(
		zap.Fields(
			zap.String("label", label),
			zap.String("app_version", "1.0"),
			zap.String("environment", "production"),
		),
	)
	if err != nil {
		fmt.Errorf("can't initialize zap logger: %v", err)
	}
	defer _logger.Sync()

	logger := Logger{
		logLevel: cfg.Service.LogLevel,
		label:    label,
		_logger:  _logger,
	}

	return &logger
}

func (l *Logger) doLog(message string, logLevel common.LogLevel) {
	isLogAllowed := logLevel >= l.logLevel
	if !isLogAllowed {
		return
	}

	switch logLevel {
	case common.DEBUG:
		l._logger.Debug(message)
	case common.INFO:
		l._logger.Info(message)
	case common.WARN:
		l._logger.Warn(message)
	case common.ERROR:
		l._logger.Error(message)
	case common.PANIC:
		l._logger.Panic(message)
	default:
		l._logger.Debug(message)
	}
}

func (l *Logger) doLogWithCtx(ctx context.Context, message string, logLevel common.LogLevel) {
	correlationID := utils.ExtractCorrelationID(ctx)

	msgWithCorrelationID := fmt.Sprintf("%s | %s", correlationID, message)

	l.doLog(msgWithCorrelationID, logLevel)
}

// normal
func (l *Logger) Debug(message string) {
	l.doLog(message, common.DEBUG)
}
func (l *Logger) Info(message string) {
	l.doLog(message, common.INFO)
}
func (l *Logger) Warn(message string) {
	l.doLog(message, common.WARN)
}
func (l *Logger) Error(message string) {
	l.doLog(message, common.ERROR)
}
func (l *Logger) Panic(message string) {
	l.doLog(message, common.PANIC)
}

// with ctx
func (l *Logger) DebugWithCtx(ctx context.Context, message string) {
	l.doLogWithCtx(ctx, message, common.DEBUG)
}
func (l *Logger) InfoWithCtx(ctx context.Context, message string) {
	l.doLogWithCtx(ctx, message, common.INFO)
}
func (l *Logger) WarnWithCtx(ctx context.Context, message string) {
	l.doLogWithCtx(ctx, message, common.WARN)
}
func (l *Logger) ErrorWithCtx(ctx context.Context, message string) {
	l.doLogWithCtx(ctx, message, common.ERROR)
}
func (l *Logger) PanicWithCtx(ctx context.Context, message string) {
	l.doLogWithCtx(ctx, message, common.PANIC)
}
