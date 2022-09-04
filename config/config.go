package config

import (
	"context"
	"errors"
	"time"

	"github.com/MrZhangjicheng/kitdemo/log"
)

var _ (Config) = (*config)(nil)

// 配置最高级接口  加载文件(一次性，以及不断监听--满足热更) 关闭 映射到不同类型的结构体
type Config interface {
	// 加载文件
	Load() error
	// 将文件中的内容进行映射
	Scan(v interface{}) error
	// 获取值 有多种类型的可能
	Value(key string) Value
	// 监听某个k,v 内容变更
	Watch() error
	// 关闭
	Close() error
}

type config struct {
	// 需要用户配置的都放在options中
	opts options
	// 进行多个文件的数据处理
	reader Reader

	watchers []Watcher
}

func New(ops ...Option) *config {
	opts := options{
		decoder:  defaultDecoder,
		resolver: defaultResolver,
	}

	for _, o := range ops {
		o(&opts)
	}

	return &config{
		opts: opts,
	}
}

func (c *config) watch(w Watcher) {
	for {
		kvs, err := w.Next()
		if errors.Is(err, context.Canceled) {
			log.Info("watcher's ctx cancel : %v", err)
			return
		}
		if err != nil {
			time.Sleep(time.Second)
			log.Errorf("failed to watch next  config: %v", err)
			continue
		}
		if err := c.reader.Merge(kvs...); err != nil {
			log.Errorf("failed to merge next config: %v", err)
			continue
		}
		if err := c.reader.Resolver(); err != nil {
			log.Errorf("failed to resolve next config: %v", err)
			continue
		}

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
		if err = c.reader.Merge(kvs...); err != nil {
			log.Errorf("failed to merge config source: %v", err)
			return err
		}
		// 创建监听者
		w, err := src.Watch()
		if err != nil {
			log.Errorf("failed to watch config source: %v", err)
			return err
		}
		c.watchers = append(c.watchers, w)
		go c.watch(w)
	}

	return nil
}

func (c *config) Close() error {
	for _, w := range c.watchers {
		if err := w.Stop(); err != nil {
			return err
		}
	}
	return nil
}
