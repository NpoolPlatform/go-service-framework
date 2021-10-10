package logger

import (
	"time"

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

func Init(level, logFile string) error {
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

	ljLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     15,
		Compress:   true,
	}
	defer ljLogger.Close()

	buildConfig := zap.NewProductionConfig()

	buildConfig.Level = zap.NewAtomicLevelAt(zapLevel)
	buildConfig.OutputPaths = []string{logFile, "stdout"}
	buildConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.StampMilli)

	_myLogger, err := buildConfig.Build()
	if err != nil {
		return xerrors.Errorf("fail to build logger: %v", err)
	}
	myLogger = _myLogger

	return nil
}

func Sugar() *zap.SugaredLogger {
	return myLogger.Sugar()
}
