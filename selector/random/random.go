package random

import (
	"context"
	"math/rand"

	"github.com/MrZhangjicheng/kitdemo/selector"
)

const (
	Name = "random"
)

var _ selector.Balancer = &Balancer{}

type Option func(o *options)

type options struct{}

type Balancer struct{}

func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, error) {
	if len(nodes) == 0 {
		return nil, selector.ErrNoAvailable
	}
	cur := rand.Intn(len(nodes))
	selectd := nodes[cur]
	// d := selectd.Pick()
	return selectd, nil

}

func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

func NewBuilder(opts ...Option) selector.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
	return &selector.DefaultBuilder{
		Balancer: &Builder{},
	}

}

type Builder struct{}

func (b *Builder) Build() selector.Balancer {
	return &Balancer{}
}
