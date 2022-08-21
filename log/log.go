package log

import (
	"context"
	"log"
)

// logger 底层的日志接口
// 日志为了排查，所以需要加一些前缀，使问题能够快速定位 比如文件名，行号
// 设置等级，可以通过等级设置，减少日志量
// 由于业务不同，需要加入的前缀不同，支持动态添加一下日志头
// 统一了日志的接入方式

var DefaultLogger = NewStdLogger(log.Writer())

type Logger interface {
	Log(level Level, keyval ...interface{}) error
}

type logger struct {
	logger    Logger
	prefix    []interface{}
	hasValuer bool
	ctx       context.Context
}

func (l *logger) Log(level Level, keyval ...interface{}) error {
	kvs := make([]interface{}, 0, len(l.prefix)+len(keyval))

	kvs = append(kvs, l.prefix...)

	if l.hasValuer {
		bindValues(l.ctx, kvs)
	}
	kvs = append(kvs, keyval...)

	if err := l.logger.Log(level, kvs...); err != nil {
		return err
	}

	return nil

}

func With(l Logger, kv ...interface{}) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{
			logger:    l,
			prefix:    kv,
			hasValuer: containsValuer(kv),
			ctx:       context.Background(),
		}
	}
	kvs := make([]interface{}, 0, len(c.prefix)+len(kv))

	kvs = append(kvs, c.prefix...)

	kvs = append(kvs, kv...)

	return &logger{
		logger:    c.logger,
		prefix:    kvs,
		hasValuer: containsValuer(kvs),
		ctx:       c.ctx,
	}
}

func WithContext(ctx context.Context, l Logger) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{
			logger: l,
			ctx:    ctx,
		}
	}
	return &logger{
		logger:    c.logger,
		prefix:    c.prefix,
		hasValuer: c.hasValuer,
		ctx:       ctx,
	}
}
