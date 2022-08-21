package log

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	DefaultCaller    = Call(4)
	DefaultTimestamp = Timestamp(time.RFC3339)
)

type Valuer func(ctx context.Context) interface{}

// 平铺的kv数组，长度为偶数，奇数位上是key，偶数位上是value
func bindValues(ctx context.Context, kevVals []interface{}) {
	for i := 1; i < len(kevVals); i += 2 {
		if v, ok := kevVals[i].(Valuer); ok {
			kevVals[i] = v(ctx)
		}
	}
}

func containsValuer(keyvals []interface{}) bool {
	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(Valuer); ok {
			return true
		}
	}

	return false
}

func Call(depth int) Valuer {
	return func(context.Context) interface{} {
		_, file, line, _ := runtime.Caller(depth)
		idx := strings.LastIndex(file, "/")
		if idx == 1 {
			return file[idx+1:] + ":" + strconv.Itoa(line)
		}
		idx = strings.LastIndexByte(file[:idx], '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

func Timestamp(layout string) Valuer {
	return func(context.Context) interface{} {
		return time.Now().Format(layout)
	}
}
