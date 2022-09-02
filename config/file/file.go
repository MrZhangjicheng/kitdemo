package file

import (
	"os"

	"github.com/MrZhangjicheng/kitdemo/config"
)

var _ config.Source = (*file)(nil)

type file struct {
	path string
}

// 本地文件需要考虑 是目录还是文件
func (f *file) Load() (kvs []*config.KeyValue, err error) {
	fi, err := os.Stat(f.path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return f.loadDir(f.path)
	}
	kv, err := f.loadFile(f.path)
	if err != nil {
		return nil, err
	}
	return []*config.KeyValue{kv}, nil

}
