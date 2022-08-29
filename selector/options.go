package selector

// 负载均衡策略
// 随机数 p2c ...
type SelectOptions struct {
}

type SelectOption func(*SelectOptions)
