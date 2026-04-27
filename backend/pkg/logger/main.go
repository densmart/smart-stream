package logger

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	layout = "2006-01-02 15:04:05"
)

var (
	infoCore  zapcore.Core
	debugCore zapcore.Core
	errorCore zapcore.Core
	fatalCore zapcore.Core
	atom      zap.AtomicLevel
)

func InitLogger() {
	atom = zap.NewAtomicLevel()
	switch viper.GetString("logger.level") {
	case "debug":
		atom.SetLevel(zap.DebugLevel)
	case "info":
		atom.SetLevel(zap.InfoLevel)
	case "error":
		atom.SetLevel(zap.ErrorLevel)
	case "fatal":
		atom.SetLevel(zap.FatalLevel)
	default:
		log.Fatalf("unknown log level: %s", viper.GetString("logger.level"))
	}

	infoCore = zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			CallerKey:  "",
			EncodeTime: zapcore.TimeEncoderOfLayout(layout),
			LineEnding: zapcore.DefaultLineEnding,
			MessageKey: "message",
			TimeKey:    "timestamp",
		}),
		zapcore.AddSync(os.Stdout), zap.InfoLevel,
	)

	debugCore = zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeTime:   zapcore.TimeEncoderOfLayout(layout),
			LineEnding:   zapcore.DefaultLineEnding,
			MessageKey:   "message",
			TimeKey:      "timestamp",
		}),
		zapcore.AddSync(os.Stdout), zap.DebugLevel,
	)

	errorCore = zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			EncodeTime:    zapcore.TimeEncoderOfLayout(layout),
			LineEnding:    zapcore.DefaultLineEnding,
			MessageKey:    "message",
			StacktraceKey: "stacktrace",
			TimeKey:       "timestamp",
		}),
		zapcore.AddSync(os.Stdout), zap.ErrorLevel,
	)

	fatalCore = zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeTime:   zapcore.TimeEncoderOfLayout(layout),
			LineEnding:   zapcore.DefaultLineEnding,
			MessageKey:   "message",
			TimeKey:      "timestamp",
		}),
		zapcore.AddSync(os.Stdout), zap.FatalLevel,
	)
}

func GetLevel() string {
	return atom.String()
}

func Infof(template string, args ...any) {
	if atom.Enabled(zap.InfoLevel) {
		zap.New(infoCore, zap.AddCaller()).Sugar().Infof(template, args...)
	}
}

func Errorf(template string, args ...any) {
	if atom.Enabled(zap.ErrorLevel) {
		zap.New(errorCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar().Errorf(template, args...)
	}
}

func Debugf(template string, args ...any) {
	if atom.Enabled(zap.DebugLevel) {
		zap.New(debugCore, zap.AddCaller(), zap.AddStacktrace(zapcore.DebugLevel)).Sugar().Debugf(template, args...)
	}
}

func Fatalf(template string, args ...any) {
	if atom.Enabled(zap.FatalLevel) {
		zap.New(fatalCore, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar().Fatalf(template, args...)
	}
}
