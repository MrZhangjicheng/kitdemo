package transport

import (
	"context"
	"net/url"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
	// 考虑是否支持多次重启
	// StartUtil(context.Context) error
}

// 该接口的作用就是为了防止用户传入的参数中没有对应的路由,然后通过server的实例进行强转
type Endpointer interface {
	Endpoint() (*url.URL, error)
}
