package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var (
	once     sync.Once
	instance *zap.Logger
)

type Config struct {
	Production bool
	LogFile    string // Path to log file (empty for no file logging)
	MaxSize    int    // Max size in MB before rotation
	MaxBackups int    // Max number of old log files
	MaxAge     int    // Max days to retain log files
	Compress   bool
}

// Init initializes the singleton logger instance
func Init(cfg Config) {
	once.Do(func() {
		var zapConfig zap.Config
		if cfg.Production {
			zapConfig = zap.NewProductionConfig()
		} else {
			zapConfig = zap.NewDevelopmentConfig()
			zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			zapConfig.OutputPaths = []string{
				"app.log",
				"stderr",
			}
		}

		cores := []zapcore.Core{}
		// Console output (stderr)
		if !cfg.Production {
			consoleEncoder := zapcore.NewConsoleEncoder(zapConfig.EncoderConfig)
			cores = append(cores, zapcore.NewCore(
				consoleEncoder,
				zapcore.Lock(os.Stderr),
				zapConfig.Level,
			))
		}

		// File output with Lumberjack rotation
		if cfg.LogFile != "" {
			fileEncoder := zapcore.NewJSONEncoder(zapConfig.EncoderConfig)
			lumberjackLogger := &lumberjack.Logger{
				Filename:   cfg.LogFile,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			}

			cores = append(cores, zapcore.NewCore(
				fileEncoder,
				zapcore.AddSync(lumberjackLogger),
				zapConfig.Level,
			))
		}

		// Combine cores
		core := zapcore.NewTee(cores...)

		// Create the logger
		instance = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

		//instance, err = zapConfig.Build()
		//if err != nil {
		//	return
		//}

		if !cfg.Production {
			instance.Info("Logger initialized in development mode")
		} else {
			instance.Info("Logger initialized in production mode")
		}

	})
}

// Get returns the singleton logger instance
func Get() *zap.Logger {
	if instance == nil {
		panic("logger not initialized")
	}
	return instance
}

// Sync flushes any buffered log entries
func Sync() error {
	if instance != nil {
		return instance.Sync()
	}
	return nil
}
