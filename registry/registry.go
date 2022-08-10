package registry

import "context"

// 注册中心
type Registar interface {
	Register(ctx context.Context) error
	Deregister(ctx context.Context) error
}

// 服务发现
type Discover interface {
	FindService(ctx context.Context)
}

// 服务的抽象,即需要在注册中心注册的信息
type ServiceIntance struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
