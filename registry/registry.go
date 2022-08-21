package registry

import "context"

// 注册中心
type Registar interface {
	Register(ctx context.Context, service *ServiceIntance) error
	Deregister(ctx context.Context, service *ServiceIntance) error
}

// 服务发现
type Discover interface {
	FindService(ctx context.Context, serviceName string) ([]*ServiceIntance, error)
}

// 服务的抽象,即需要在注册中心注册的信息  k-v
type ServiceIntance struct {
	ID        string   `json:"id"`   // 服务实例唯一id
	Name      string   `json:"name"` // 服务名称
	Version   string   `json:"version"`
	Endpoints []string `json:"endpoints"` // 服务地址   “http://127.0.0.1:8000”,"grpc://127.0.0.1:9000"
}
