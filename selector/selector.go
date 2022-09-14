package selector

import "context"

// 路由与负载均衡的高级接口
/*
	三种方向
		节点权重计算
		服务路由过滤策略
		负载均衡算法
*/
type Selector interface {
	Select(ctx context.Context, opts ...SelectOption) (selected Node, err error)
}

// 构建一个选择器
type Builder interface {
	Build() Selector
}

// 从注册中心拿到的节点列表
type Node interface {
	Scheme() string

	Address() string

	ServiceName() string

	Version() string
}
