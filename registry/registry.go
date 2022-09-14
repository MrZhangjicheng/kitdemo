package registry

import "context"

// 该模块是注册中心的最高级接口，所有要集成的注册中心要实现该接口
type Registar interface {
	Register(ctx context.Context, service *ServiceIntance) error
	Deregister(ctx context.Context, service *ServiceIntance) error
}

// 服务发现
type Discover interface {
	FindService(ctx context.Context, serviceName string) ([]*ServiceIntance, error)
	// 创建一个观察者
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

type Watcher interface {
	// 第一观看且服务列表不为空
	// 任何实例列表更改
	Next() ([]*ServiceIntance, error)

	Stop() error
}

// 服务的抽象,即需要在注册中心注册的信息  k-v
type ServiceIntance struct {
	ID        string            `json:"id"`   // 服务实例唯一id
	Name      string            `json:"name"` // 服务名称
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []string          `json:"endpoints"` // 服务地址   “http://127.0.0.1:8000”,"grpc://127.0.0.1:9000"
}
