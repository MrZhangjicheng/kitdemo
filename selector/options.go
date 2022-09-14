package selector

// 节点筛选  通过版本号
type SelectOptions struct {
	NodeFilters []NodeFilter
}

type SelectOption func(*SelectOptions)

func WithNodeFilter(fn ...NodeFilter) SelectOption {
	return func(opts *SelectOptions) {
		opts.NodeFilters = fn
	}
}
