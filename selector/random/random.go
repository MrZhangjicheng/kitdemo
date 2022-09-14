package random

import (
	"github.com/MrZhangjicheng/kitdemo/selector"
)

const (
	Name = "random"
)

type Option func(o *options)

type options struct{}

func NewBuilder(opts ...Option) selector.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
}
