package config

import (
	"fmt"
	"strings"

	"github.com/MrZhangjicheng/kitdemo/encoding"
)

type Decoder func(*KeyValue, map[string]interface{}) error

type Resolver func(map[string]interface{}) error

type Option func(*options)

// 需要用户设置  文件 编解码方式
type options struct {
	sources  []Source
	decoder  Decoder
	resolver Resolver
}

func WithSources(sources ...Source) Option {
	return func(o *options) {
		o.sources = sources
	}

}

func WithDecoder(d Decoder) Option {
	return func(o *options) {
		o.decoder = d
	}
}

func WithResolver(r Resolver) Option {
	return func(o *options) {
		o.resolver = r
	}
}

func defaultDecoder(src *KeyValue, target map[string]interface{}) error {
	// 针对未支持的类型
	if src.Format == "" {
		keys := strings.Split(src.Key, ".")
		for i, k := range keys {
			if i == len(keys)-1 {
				target[k] = src.Value
			} else {
				sub := make(map[string]interface{})
				target[k] = sub
				target = sub
			}
		}
		return nil
	}
	if codec := encoding.GetCodec(src.Format); codec != nil {
		return codec.Unmarshal(src.Value, &target)
	}

	return fmt.Errorf("unsupporten key: %s format: %s", src.Key, src.Format)
}

func defaultResolver(input map[string]interface{}) error {
	return nil

}
