package discovery

import (
	"time"

	"github.com/MrZhangjicheng/kitdemo/registry"

	"google.golang.org/grpc/resolver"
)

const name = "discovery"

type Option func(o *builder)

func WithTimeout(timeout time.Duration) Option {
	return func(o *builder) {
		o.timeout = timeout
	}
}

func WithInsecure(insecure bool) Option {
	return func(o *builder) {
		o.insecure = insecure
	}
}

func DisableDebugLog() Option {
	return func(o *builder) {
		o.debugLogDisabled = true
	}
}

// 这个实现grpc内部接口，来实现负载均衡
type builder struct {
	discover         registry.Discover
	timeout          time.Duration
	insecure         bool
	debugLogDisabled bool
}

func NewBuilder(d registry.Discover, opts ...Option) resolver.Builder {
	b := &builder{
		discover:         d,
		timeout:          time.Second * 10,
		insecure:         false,
		debugLogDisabled: false,
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// watchRes := &struct {
	// 	err error
	// 	w   registry.Watcher
	// }{}
	// done := make(chan struct{}, 1)
	// ctx, cancel := context.WithCancel(context.Background())
	// go func() {

	// }
	r := &discoveryResolver{}

	return r, nil

}

func (*builder) Scheme() string {
	return name
}
