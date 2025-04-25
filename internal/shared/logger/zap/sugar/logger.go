package logger

import (
	"strings"

	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"go.uber.org/zap"
)

type LogLevel = logger.LogLevel

const (
	LevelDebug   = logger.LevelDebug
	LevelInfo    = logger.LevelInfo
	LevelWarning = logger.LevelWarning
	LevelError   = logger.LevelError
)

type ZapSugarLogger struct {
	logger      *zap.SugaredLogger
	minLogLevel LogLevel
}

func New(minLogLevel LogLevel) *ZapSugarLogger {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zslogger := &ZapSugarLogger{
		logger:      log.Sugar(),
		minLogLevel: minLogLevel,
	}
	zslogger.Info(
		"ZapSugarLogger created",
		"level", logger.GetLevelName(minLogLevel),
	)
	return zslogger
}

func (l *ZapSugarLogger) Debug(message string, keysAndValues ...interface{}) {
	if LevelDebug >= l.minLogLevel {
		// l.logger.Debugw(message, keysAndValues...)
		l.logger.Infow(message, keysAndValues...)
	}
}

func (l *ZapSugarLogger) Info(message string, keysAndValues ...interface{}) {
	if LevelInfo >= l.minLogLevel {
		l.logger.Infow(message, keysAndValues...)
	}
}

func (l *ZapSugarLogger) Warning(message string, keysAndValues ...interface{}) {
	if LevelWarning >= l.minLogLevel {
		// l.logger.Warnw(message, keysAndValues...)
		l.logger.Infow(message, keysAndValues...)
	}
}

func (l *ZapSugarLogger) Error(message string, err error, keysAndValues ...interface{}) {
	if LevelError >= l.minLogLevel {
		addKeysAndValues := append([]interface{}{"reason", err.Error()}, keysAndValues...)
		// l.logger.Errorw(message, addKeysAndValues...)
		l.logger.Infow(message, addKeysAndValues...)
	}
}

func (l *ZapSugarLogger) Close() {
	l.Info(
		"ZapSugarLogger closed",
		"level", logger.GetLevelName(l.minLogLevel),
	)
	if err := l.logger.Sync(); err != nil {
		if !strings.Contains(err.Error(), "The handle is invalid") {
			panic(err)
		}
	}
}
