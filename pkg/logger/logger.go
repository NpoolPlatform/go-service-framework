package logger

import (
	"fmt"
	"os"

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

var _logger *zap.Logger

// https://github.com/uber-go/zap/blob/master/FAQ.md
func Init(level, logFile string, opts ...zap.Option) error {
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
		return fmt.Errorf("unknow log level %s", level)
	}

	fileLog := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     15,
		Compress:   true,
	})

	encode := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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
	})

	multiSyncLog := zapcore.NewMultiWriteSyncer(
		fileLog,
		zapcore.AddSync(os.Stdout),
	)

	core := zapcore.NewCore(encode, multiSyncLog, zapLevel)
	opts = append(opts, zap.AddCaller())
	_logger = zap.New(core).WithOptions(opts...)
	return nil
}

func Sugar() *zap.SugaredLogger {
	return _logger.Sugar()
}

func Sync() error {
	return _logger.Sync()
}
