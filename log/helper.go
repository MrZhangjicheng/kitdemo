package log

import (
	"context"
	"fmt"
	"os"
)

var DefaultMessageKey = "msg"

type Option func(*Helper)

// 日志的高级接口，建议使用
type Helper struct {
	logger Logger
	msgKey string
}

func WithMessageKey(k string) Option {
	return func(opts *Helper) {
		opts.msgKey = k
	}
}

func NewHelper(logger Logger, opts ...Option) *Helper {
	options := &Helper{
		msgKey: DefaultMessageKey,
		logger: logger,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func (h *Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{
		msgKey: h.msgKey,
		logger: WithContext(ctx, h.logger),
	}
}

func (h *Helper) Log(level Level, keyvals ...interface{}) {
	_ = h.logger.Log(level, keyvals...)
}

// Debug logs a message at debug level.
func (h *Helper) Debug(a ...interface{}) {
	_ = h.logger.Log(LevelDebug, h.msgKey, fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func (h *Helper) Debugf(format string, a ...interface{}) {
	_ = h.logger.Log(LevelDebug, h.msgKey, fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *Helper) Debugw(keyvals ...interface{}) {
	_ = h.logger.Log(LevelDebug, keyvals...)
}

// Info logs a message at info level.
func (h *Helper) Info(a ...interface{}) {
	_ = h.logger.Log(LevelInfo, h.msgKey, fmt.Sprint(a...))
}

// Infof logs a message at info level.
func (h *Helper) Infof(format string, a ...interface{}) {
	_ = h.logger.Log(LevelInfo, h.msgKey, fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func (h *Helper) Infow(keyvals ...interface{}) {
	_ = h.logger.Log(LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func (h *Helper) Warn(a ...interface{}) {
	_ = h.logger.Log(LevelWarn, h.msgKey, fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *Helper) Warnf(format string, a ...interface{}) {
	_ = h.logger.Log(LevelWarn, h.msgKey, fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (h *Helper) Warnw(keyvals ...interface{}) {
	_ = h.logger.Log(LevelWarn, keyvals...)
}

// Error logs a message at error level.
func (h *Helper) Error(a ...interface{}) {
	_ = h.logger.Log(LevelError, h.msgKey, fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func (h *Helper) Errorf(format string, a ...interface{}) {
	_ = h.logger.Log(LevelError, h.msgKey, fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func (h *Helper) Errorw(keyvals ...interface{}) {
	_ = h.logger.Log(LevelError, keyvals...)
}

// Fatal logs a message at fatal level.
func (h *Helper) Fatal(a ...interface{}) {
	_ = h.logger.Log(LevelFatal, h.msgKey, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func (h *Helper) Fatalf(format string, a ...interface{}) {
	_ = h.logger.Log(LevelFatal, h.msgKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func (h *Helper) Fatalw(keyvals ...interface{}) {
	_ = h.logger.Log(LevelFatal, keyvals...)
	os.Exit(1)
}
