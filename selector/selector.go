package selector

import (
	"context"
	"errors"
)

var ErrNoAvailable = errors.New("no_available_node")

// 路由与负载均衡的高级接口
/*
	三种方向
		节点权重计算
		服务路由过滤策略
		负载均衡算法
	意义：
		1. 通过合适的算法选择真正的服务节点
		2. 监听客户端的服务列表,并及时更新
*/

type Selector interface {
	Rebalancer
	Select(ctx context.Context, opts ...SelectOption) (selected Node, err error)
}

// 节点平衡器
type Rebalancer interface {
	// 发生改变应用到所有节点
	Apply(nodes []Node)
}

// 构建一个选择器
type Builder interface {
	Build() Selector
}

// 从注册中心拿到的节点列表抽象
type Node interface {
	Scheme() string

	Address() string

	ServiceName() string

	Version() string
}
