package logger

import (
	"testing"
)

func TestInit(t *testing.T) {
	err := logger.Init(logger.DebugLevel, "./debug.log")
	if err != nil {
		t.Errorf("Fail to init logger with debug level: %v", err)
	}
}
