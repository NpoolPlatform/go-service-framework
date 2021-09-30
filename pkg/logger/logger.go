package logger

import (
	"golang.org/x/xerrors"

	zap "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DebugLevel   = "debug"
	InfoLevel    = "info"
	WarningLevel = "warning"
	ErrorLevel   = "error"
)

var myLogger *zap.Logger

func Init(level string, logFile string) error {
	var zapLevel zapcore.Level
	switch level {
	case DebugLevel:
		zapLevel = zap.DebugLevel
	case InfoLevel:
		zapLevel = zap.InfoLevel
	case WarningLevel:
		zapLevel = zap.WarnLevel
	case ErrorLevel:
		zapLevel = zap.ErrorLevel
	default:
		return xerrors.Errorf("unknow log level %s", level)
	}

	encCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	ljLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     15,
		Compress:   true,
	}
	defer ljLogger.Close()

	myLogger = zap.New(
		zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg),
			zapcore.AddSync(ljLogger),
			zapLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return nil
}

func Sugar() *zap.SugaredLogger {
	return myLogger.Named("sugar").Sugar()
}
