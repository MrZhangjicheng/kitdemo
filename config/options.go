package config

import (
	"fmt"
	"log"
	"regexp"
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

func WithLogger(l log.Logger) Option {
	return func(o *options) {}
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

// 逻辑不清，应该是针对不同种类型进行解析
func defaultResolver(input map[string]interface{}) error {
	mapper := func(name string) string {
		args := strings.SplitN(strings.TrimSpace(name), ":", 2)
		if v, has := readValue(input, args[0]); has {
			s, _ := v.String()
			return s
		} else if len(args) > 1 {
			return args[1]
		}

		return ""
	}
	var resolve func(map[string]interface{}) error
	resolve = func(sub map[string]interface{}) error {
		for k, v := range sub {
			switch vt := v.(type) {
			case string:
				sub[k] = expand(vt, mapper)
			case map[string]interface{}:
				if err := resolve(vt); err != nil {
					return err
				}
			case []interface{}:
				for i, iface := range vt {
					switch it := iface.(type) {
					case string:
						vt[i] = expand(it, mapper)
					case map[string]interface{}:
						if err := resolve(it); err != nil {
							return err
						}
					}
				}
				sub[k] = vt
			}
		}
		return nil
	}
	return resolve(input)

}

// 变量替换 但具体实现不清楚
func expand(s string, mapping func(string) string) string {
	r := regexp.MustCompile(`\${(.*?)}`)
	re := r.FindAllStringSubmatch(s, -1)
	for _, i := range re {
		if len(i) == 2 {
			s = strings.ReplaceAll(s, i[0], mapping(i[1]))
		}
	}

	return s
}
