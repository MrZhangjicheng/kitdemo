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
	// 创建一个观察者
	Watch() (Watcher, error)
}

// 观察者 进行热更新
type Watcher interface {
	Next() ([]*KeyValue, error)

	Stop() error
}
