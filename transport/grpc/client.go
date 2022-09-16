package grpc

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/MrZhangjicheng/kitdemo/registry"
	"github.com/MrZhangjicheng/kitdemo/selector"
	"github.com/MrZhangjicheng/kitdemo/transport/grpc/resolver/discovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
)

type ClientOption func(*clientOptions)

func WithEndPoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

func Withdiscovery(discovery registry.Discover) ClientOption {
	return func(o *clientOptions) {
		o.discovery = discovery
	}
}

// 拦截器
func WithUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.ints = in
	}
}

// 一些配置
func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.grpcOpts = opts
	}
}

func WithLogger(log log.Logger) ClientOption {
	return func(o *clientOptions) {}
}

// 需要转化成 grpc
func WithNodeFilter(filters ...selector.NodeFilter) ClientOption {
	return func(o *clientOptions) {
		o.filters = filters
	}
}

type clientOptions struct {
	// 客户端基本配置
	endpoint string // 两种模式 一种是直接ip+port  另一种采用服务发现 "discovery//<author>/servicename"
	tlsConf  *tls.Config
	// 超时控制
	timeout time.Duration
	// 客户端拦截器
	ints []grpc.UnaryClientInterceptor
	// grpc 连接参数配置相关
	grpcOpts  []grpc.DialOption
	discovery registry.Discover

	balancerName string
	// 过滤器 中间件
	filters []selector.NodeFilter
}

// Dial 返回 grpc 的 连接
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout: 2000 * time.Microsecond,
	}
	for _, o := range opts {
		o(&options)
	}
	// 将自定义的中间件,过滤器 装华为 grpc的拦截器
	ints := []grpc.UnaryClientInterceptor{}
	if len(options.ints) > 0 {
		ints = append(ints, options.ints...)
	}
	grpcOpts := []grpc.DialOption{
		// 负载均衡相关配置
		// grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalacncingConfig":[{"%s":{}}]}`, options.balancerName)),

		grpc.WithChainUnaryInterceptor(ints...),
	}
	if options.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers(
			discovery.NewBuilder(
				options.discovery,
				discovery.WithInsecure(insecure),
			)))
	}
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}
	if options.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(options.tlsConf)))
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}
	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}
