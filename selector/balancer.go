package selector

import "context"

type Balancer interface {
	Pick(ctx context.Context, nodes []WeightedNode) (selected WeightedNode, err error)
}

type BalancerBuilder interface {
	Build() Balancer
}

// 加权节点 接口
type WeightedNode interface {
	Node
	// 返回原始节点
	Raw() Node
	// 运行时的权重计算
	Weight() float64
}

type WeightedNodeBuilder interface {
	Build(Node) WeightedNode
}
