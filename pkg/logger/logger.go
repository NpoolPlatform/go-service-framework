package logger

import (
	"golang.org/x/xerrors"
	"net/url"

	zap "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
	lj "gopkg.in/natefinch/lumberjack.v2"
)

const (
	DebugLevel   = "debug"
	InfoLevel    = "info"
	WarningLevel = "warning"
	ErrorLevel   = "error"
)

var myLogger zap.Logger

func Init(level string, logFile string) error {
	var zapLevel zapcore.Level
	switch level {
	case DebugLevel:
		zapLevel = zap.DebugLevel
	case InfoLevel:
		zapLevel = zap.InfoLevel
	case WarningLevel:
		zapLevel = zap.WarningLevel
	case ErrorLevel:
		zapLevel = zap.ErrorLevel
	default:
		return xerrors.Errorf("unknow log level %s", level)
	}

	encCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallKey:        "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	ljLogger := lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     15,
		Compress:   true,
	}

	zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
		return &ljLogger, nil
	})

	myLogger, err := loggerConfig.Build()
	if err != nil {
		return xerror.Errorf("fail to build logger config: %v", err)
	}

	zap.ReplaceGlobals(myLogger)
	return nil
}
