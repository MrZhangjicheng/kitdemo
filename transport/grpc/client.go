package grpc

import (
	"context"
	"crypto/tls"

	"github.com/MrZhangjicheng/kitdemo/registry"
	"github.com/MrZhangjicheng/kitdemo/transport/grpc/resolver/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type clientOptions struct {
	// 客户端基本配置
	endpoint  string
	tlsConf   *tls.Config
	ints      []grpc.UnaryClientInterceptor
	grpcOpts  []grpc.DialOption
	discovery registry.Discover
}

func dial(ctx context.Context, insecure bool, opts ...clientOptions) (*grpc.ClientConn, error) {
	options := clientOptions{}
	ints := []grpc.UnaryClientInterceptor{}
	if len(options.ints) > 0 {
		ints = append(ints, options.ints...)
	}
	grpcOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(ints...),
	}
	if options.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers(
			discovery.NewBuilder(
				options.discovery,
				// discovery.WithInsecure(insecure),
			)))
	}
	if options.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(options.tlsConf)))
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}
	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}
