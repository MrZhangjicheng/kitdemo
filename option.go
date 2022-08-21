package kitdemo

import (
	"context"

	"github.com/MrZhangjicheng/kitdemo/log"

	"net/url"
	"os"

	"github.com/MrZhangjicheng/kitdemo/registry"
	"github.com/MrZhangjicheng/kitdemo/transport"
)

// 用户必须要传的
// 服务以及对应的注册中心
// 服务名称以及对应的 访问地址
// 服务接受的信号

type Option func(o *options)

type options struct {
	//注册服务实例的基本信息
	id        string
	name      string
	version   string
	endpoints []*url.URL
	// 设计到生命周期的管理
	ctx     context.Context
	singles []os.Signal

	// 注册中心相关
	registar registry.Registar

	// 服务相关 http grpc
	Servers []transport.Server

	// 日志
	logger log.Logger
}

//  options 不对外暴露，通过 Option来设置参数
func ID(id string) Option {
	return func(o *options) { o.id = id }
}
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

func Version(version string) Option {
	return func(o *options) { o.version = version }
}

func Endpoints(endpoints ...*url.URL) Option {
	return func(o *options) { o.endpoints = endpoints }
}
func Registar(registar registry.Registar) Option {
	return func(o *options) { o.registar = registar }
}
func Server(srv ...transport.Server) Option {
	return func(o *options) { o.Servers = srv }
}
func Singles(s ...os.Signal) Option {
	return func(o *options) { o.singles = s }
}

func Logger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}
