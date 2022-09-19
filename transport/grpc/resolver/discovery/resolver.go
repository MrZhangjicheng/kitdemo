package discovery

import (
	"google.golang.org/grpc/resolver"
)

// 实现grpc内部接口，达成客户端负载均衡
type discoveryResolver struct {
}

func (r *discoveryResolver) ResolveNow(options resolver.ResolveNowOptions) {}

func (r *discoveryResolver) Close() {}
