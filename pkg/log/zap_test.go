package log

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestInit(t *testing.T) {
	Init(zapcore.ErrorLevel)
	logger := Logger()
	defer logger.Sync()

	Init(zapcore.ErrorLevel)

	if logger != Logger() {
		t.Fatal("logger got reinitialized, expected not to do so")
	}
}
