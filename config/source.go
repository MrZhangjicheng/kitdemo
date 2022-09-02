package config

// 文件内容实际存储的结构
type KeyValue struct {
	Key    string
	Value  []byte
	Format string
}

type Source interface {
	// 数据源本身进行加载处理
	Load() ([]*KeyValue, error)
}
