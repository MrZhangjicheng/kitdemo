package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"

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

type clientOptions struct {
	// 客户端基本配置
	endpoint string // 两种模式 一种是直接ip+port  另一种采用服务发现 "discovery//<author>/servicename"
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

func Dial(ctx context.Context, opts ...Option) (*grpc.ClientConn, error) {
	options := clientOptions{}
	grpcOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	for _, opt := range opts {
		opt(&options)
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}
	// 判断 endpoint 传入的是 服务的地址还是注册中心的地址

	target, err := parseTarget(options.endpoint)
	if err != nil {
		return nil, err
	}
	if target.Scheme == "grpc" && target.Authority != "" {
		return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
	}

	if options.discovery != nil {
		watcher, err := options.discovery.Watch(context.Background(), target.Endpoint)
		if err != nil {
			panic("服务发现错误")
		}
		srvs, _ := watcher.Next()
		// 同时开启http服务和grpc服务选择对服务
		srv := srvs[0]
		for _, k := range srv.Endpoints {
			if strings.Contains(k, "grpc") {
				options.endpoint = k[7:]
			}
		}
	}

	fmt.Println(options.endpoint)
	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)

}

type Target struct {
	Scheme    string
	Authority string
	Endpoint  string
}

func parseTarget(endpoint string) (*Target, error) {
	if !strings.Contains(endpoint, "://") {
		endpoint = "grpc://" + endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	target := &Target{Scheme: u.Scheme, Authority: u.Host}
	if len(u.Path) > 1 {
		target.Endpoint = u.Path[1:]
	}

	return target, nil
}
