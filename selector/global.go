package selector

// 全局selector 构建器
var globalSelector Builder

func GlobalSelector() Builder {
	return globalSelector
}

func SetGlobalSelector(builder Builder) {
	globalSelector = builder
}
