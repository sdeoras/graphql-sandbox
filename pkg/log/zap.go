package log

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var once sync.Once

// Logger retrieves the singleton value to
// the caller.
func Logger() *zap.Logger {
	return logger
}

// Init initialized zap logger.
// Initialization can happen only once.
func Init(logLevel zapcore.Level) {
	once.Do(func() {
		var level zap.LevelEnablerFunc
		level = func(lvl zapcore.Level) bool {
			return lvl >= logLevel
		}

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		encoder := zapcore.NewJSONEncoder(encoderCfg)
		consoleErrors := zapcore.Lock(os.Stderr)
		core := zapcore.NewCore(encoder, consoleErrors, level)
		logger = zap.New(core)
	})
}
