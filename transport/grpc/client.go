package grpc

import (
	"context"
	"crypto/tls"

	"github.com/MrZhangjicheng/kitdemo/registry"
	"google.golang.org/grpc"
)

type Option func(*clientOptions)

func WithEndPoint(endpoint string) Option {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

func Withdiscovery(discovery registry.Discover) Option {
	return func(o *clientOptions) {
		o.discovery = discovery
	}
}

func WithserviceName(discovery registry.Discover) Option {
	return func(o *clientOptions) {
		o.discovery = discovery
	}
}

type clientOptions struct {
	serviceName string
	// 客户端基本配置
	endpoint string
	tlsConf  *tls.Config
	// ints     []grpc.UnaryClientInterceptor
	grpcOpts  []grpc.DialOption
	discovery registry.Discover
}

// func dial(ctx context.Context, insecure bool, opts ...clientOptions) (*grpc.ClientConn, error) {
// 	options := clientOptions{}
// 	ints := []grpc.UnaryClientInterceptor{}
// 	if len(options.ints) > 0 {
// 		ints = append(ints, options.ints...)
// 	}
// 	grpcOpts := []grpc.DialOption{
// 		grpc.WithChainUnaryInterceptor(ints...),
// 	}
// 	if options.discovery != nil {
// 		grpcOpts = append(grpcOpts, grpc.WithResolvers(
// 			discovery.NewBuilder(
// 				options.discovery,
// 				// discovery.WithInsecure(insecure),
// 			)))
// 	}
// 	if options.tlsConf != nil {
// 		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(options.tlsConf)))
// 	}
// 	if len(options.grpcOpts) > 0 {
// 		grpcOpts = append(grpcOpts, options.grpcOpts...)
// 	}
// 	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
// }

func dial(ctx context.Context, opts ...Option) (*grpc.ClientConn, error) {
	options := clientOptions{}
	grpcOpts := []grpc.DialOption{}
	for _, opt := range opts {
		opt(&options)
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}
	if options.endpoint == "" && options.discovery != nil {
		watcher, err := options.discovery.Watch(context.Background(), options.serviceName)
		if err != nil {
			panic("服务发现错误")
		}
		srvs, _ := watcher.Next()
		// 任意选择一个服务实例 负载均衡
		options.endpoint = srvs[0].Endpoints[0][7:]

	}
	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)

}
