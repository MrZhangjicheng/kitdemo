package config

import "github.com/MrZhangjicheng/kitdemo/log"

var _ (Config) = (*config)(nil)

// 配置最高级接口
type Config interface {
	// 加载数据
	Load() error
}

type config struct {
	// 需要用户配置的都放在options中
	opts options
}

func New(ops ...Option) *config {
	opts := options{}

	for _, o := range ops {
		o(&opts)
	}

	return &config{
		opts: opts,
	}
}

func (c *config) Load() error {
	for _, src := range c.opts.sources {
		kvs, err := src.Load()
		if err != nil {
			return err
		}
		for _, v := range kvs {
			log.Debugf("config loaded: %s format: %s", v.Key, v.Format)
		}
	}

	return nil
}
