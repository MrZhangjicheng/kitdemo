package selector

import "context"

var (
	_ Builder  = (*DefaultBuilder)(nil)
	_ Selector = (*Default)(nil)
)

type Default struct {
}

func (d *Default) Select(ctx context.Context, opts ...SelectOption) (selected Node, err error) {

	return
}

type DefaultBuilder struct {
}

func (db *DefaultBuilder) Build() Selector {
	return &Default{}
}
