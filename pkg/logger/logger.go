package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	TimeFormat string
	TimeZone string
	TimeInterval time.Duration
	LogLevel string
	Output *os.File
	SkipPaths []string
}

var DefaultLoggerConfig = LoggerConfig{
	TimeFormat: "2006-01-02 15:04:05",
	TimeZone: "UTC",
	TimeInterval: 500 * time.Millisecond,
	LogLevel: "info",
	Output: os.Stdout,
	SkipPaths: []string{},
}

type Logger struct {
	*zap.Logger
	config LoggerConfig
}

func NewLogger(config ...LoggerConfig) *Logger {

	var cfg = DefaultLoggerConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.TimeFormat == "" {
		cfg.TimeFormat = DefaultLoggerConfig.TimeFormat
	}
	if cfg.TimeZone == "" {
		cfg.TimeZone = DefaultLoggerConfig.TimeZone
	}
	if cfg.TimeInterval == 0 {
		cfg.TimeInterval = DefaultLoggerConfig.TimeInterval
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = DefaultLoggerConfig.LogLevel
	}
	if cfg.Output == nil {
		cfg.Output = DefaultLoggerConfig.Output
	}

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		level = zapcore.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(cfg.Output),
		level,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Logger{
		Logger: logger,
		config: cfg,
	}

}

func (l *Logger) WithField(key string, value interface{}) *zap.Logger {
	return l.Logger.With(zap.Any(key, value))
}

func (l *Logger) WithFields(fields map[string]interface{}) *zap.Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return l.Logger.With(zapFields...)
}

func (l *Logger) Close() error {
	return l.Logger.Sync()
}