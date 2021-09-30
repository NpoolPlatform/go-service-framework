package logger

import (
	"testing"
)

func TestInit(t *testing.T) {
	logFile := "stderr"

	err := Init(DebugLevel, logFile)
	if err != nil {
		t.Errorf("Fail to init logger with debug level: %v", err)
	}

	sugar := Sugar()
	if sugar == nil {
		t.Errorf("Sugar is not initialized")
	}

	if sugar != nil {
		sugar.Infow("test for logger infow",
			"file", logFile)
		sugar.Infof("test for logger infof %v",
			logFile)
	}
}
