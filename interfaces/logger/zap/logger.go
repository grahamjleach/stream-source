package zap

import (
	"errors"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	errorUnsupportedLogLevel = errors.New("unsupported log level")
	logLevels                = map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
	}
)

type Logger struct {
	*zap.Logger
}

type LoggerConfig struct {
	LogLevel      string
	CaptureStdLog bool
	Context       map[string]interface{}
}

func NewJSONLogger(config *LoggerConfig) (logger *Logger, err error) {
	zapLevel, found := logLevels[config.LogLevel]
	if !found {
		err = errorUnsupportedLogLevel
		return
	}

	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		DisableCaller:    true,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    config.Context,
	}

	zapLogger, err := zapConfig.Build()
	if err != nil {
		return
	}

	logger = &Logger{
		Logger: zapLogger,
	}

	if config.CaptureStdLog {
		log.SetOutput(logger)
	}
	return
}

func (l *Logger) Debug(message string, keyvals map[string]string) {
	l.Logger.Debug(message, l.fields(keyvals)...)
}

func (l *Logger) Info(message string, keyvals map[string]string) {
	l.Logger.Info(message, l.fields(keyvals)...)
}

func (l *Logger) Warn(message string, keyvals map[string]string) {
	l.Logger.Warn(message, l.fields(keyvals)...)
}

func (l *Logger) Error(message string, keyvals map[string]string) {
	l.Logger.Error(message, l.fields(keyvals)...)
}

func (l *Logger) Fatal(message string, keyvals map[string]string) {
	l.Logger.Fatal(message, l.fields(keyvals)...)
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.Logger.Info(string(p))
	return
}

func (l *Logger) fields(keyvals map[string]string) []zapcore.Field {
	fields := make([]zapcore.Field, 0, len(keyvals))

	for k, v := range keyvals {
		fields = append(fields, zap.String(k, v))
	}

	return fields
}
