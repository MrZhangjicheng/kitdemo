package discovery

import (
	"github.com/MrZhangjicheng/kitdemo/registry"

	"google.golang.org/grpc/resolver"
)

const name = "discovery"

type builder struct {
	discover registry.Discover
}

func NewBuilder(d registry.Discover) resolver.Builder {
	b := &builder{
		discover: d,
	}

	return b
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return nil, nil

}

func (*builder) Scheme() string {
	return name
}
