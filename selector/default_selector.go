package selector

import (
	"context"
	"sync/atomic"
)

var (
	_ Builder  = (*DefaultBuilder)(nil)
	_ Selector = (*Default)(nil)
)

/*
	作用
		1. 给服务列表的节点进行计算加权值    --- 加权节点
		2. 选择合适的服务节点             --- 平衡器
		3. 将服务列表的修改及时同步       --- 平衡器
*/
type Default struct {
	NodeBuilder WeightedNodeBuilder
	Balancer    Balancer
	nodes       atomic.Value
}

// 该函数就是从本地服务列表中，通过节点过滤器进行过滤，获得最终真正的调用节点
func (d *Default) Select(ctx context.Context, opts ...SelectOption) (selected Node, err error) {
	var (
		options    SelectOptions
		candidates []WeightedNode
	)
	nodes, ok := d.nodes.Load().([]WeightedNode)
	if !ok {
		return nil, ErrNoAvailable
	}
	for _, o := range opts {
		o(&options)
	}
	if len(options.NodeFilters) > 0 {
		newNodes := make([]Node, len(nodes))
		for i, wc := range nodes {
			newNodes[i] = wc
		}
		for _, filter := range options.NodeFilters {
			newNodes = filter(ctx, newNodes)
		}
		candidates = make([]WeightedNode, len(newNodes))
		for i, n := range newNodes {
			candidates[i] = n.(WeightedNode)
		}
	} else {
		candidates = nodes
	}
	if len(candidates) == 0 {
		return nil, ErrNoAvailable
	}
	wn, err := d.Balancer.Pick(ctx, candidates)
	if err != nil {
		return nil, err
	}
	return wn.Raw(), nil

	return
}

// 如果服务列表更新，直接替换本地cache
func (d *Default) Apply(nodes []Node) {
	WeightedNodes := make([]WeightedNode, 0, len(nodes))
	for _, n := range nodes {
		WeightedNodes = append(WeightedNodes, d.NodeBuilder.Build(n))
	}
	d.nodes.Store(WeightedNodes)

}

type DefaultBuilder struct {
	Node     WeightedNodeBuilder
	Balancer BalancerBuilder
}

func (db *DefaultBuilder) Build() Selector {
	return &Default{
		NodeBuilder: db.Node,
		Balancer:    db.Balancer.Build(),
	}
}
