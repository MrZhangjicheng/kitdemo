package selector

import (
	"strconv"

	"github.com/MrZhangjicheng/kitdemo/registry"
)

// 从注册中心获得的节点 默认抽象
type DefaultNode struct {
	scheme   string
	addr     string
	version  string
	name     string
	metadata map[string]string
	// 负载均衡下使用
	weight *int64
}

func (n *DefaultNode) Scheme() string {
	return n.scheme
}

func (n *DefaultNode) Address() string {
	return n.addr
}

func (n *DefaultNode) ServiceName() string {
	return n.name
}

func (n *DefaultNode) Version() string {
	return n.version
}

func NewNode(schema, addr string, ins *registry.ServiceIntance) Node {
	n := &DefaultNode{
		scheme: schema,
		addr:   addr,
	}
	if ins != nil {
		n.name = ins.Name
		n.version = ins.Version
		n.metadata = ins.Metadata
		if str, ok := ins.Metadata["weight"]; ok {
			if weight, err := strconv.ParseInt(str, 10, 64); err == nil {
				n.weight = &weight
			}
		}

	}
	return n
}
