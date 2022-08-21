package log

import (
	"context"
	"fmt"
	"os"
	"sync"
)

var global = &loggerAppliance{}

type loggerAppliance struct {
	lock sync.Mutex
	Logger
}

func init() {
	global.SetLogger(DefaultLogger)
}

func (a *loggerAppliance) SetLogger(in Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
}

func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

func GetLogger() Logger {
	return global
}

func Context(ctx context.Context) *Helper {
	return NewHelper(WithContext(ctx, global.Logger))
}

func Log(level Level, keyvals ...interface{}) {
	_ = global.Log(level, keyvals...)
}

func Debug(a ...interface{}) {
	_ = global.Log(LevelDebug, DefaultMessageKey, fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...interface{}) {
	_ = global.Log(LevelDebug, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func Debugw(keyvals ...interface{}) {
	_ = global.Log(LevelDebug, keyvals...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	_ = global.Log(LevelInfo, DefaultMessageKey, fmt.Sprint(a...))
}

// Infof logs a message at info level.
func Infof(format string, a ...interface{}) {
	_ = global.Log(LevelInfo, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func Infow(keyvals ...interface{}) {
	_ = global.Log(LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	_ = global.Log(LevelWarn, DefaultMessageKey, fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...interface{}) {
	_ = global.Log(LevelWarn, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func Warnw(keyvals ...interface{}) {
	_ = global.Log(LevelWarn, keyvals...)
}

// Error logs a message at error level.
func Error(a ...interface{}) {
	_ = global.Log(LevelError, DefaultMessageKey, fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func Errorf(format string, a ...interface{}) {
	_ = global.Log(LevelError, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func Errorw(keyvals ...interface{}) {
	_ = global.Log(LevelError, keyvals...)
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	_ = global.Log(LevelFatal, DefaultMessageKey, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...interface{}) {
	_ = global.Log(LevelFatal, DefaultMessageKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func Fatalw(keyvals ...interface{}) {
	_ = global.Log(LevelFatal, keyvals...)
	os.Exit(1)
}