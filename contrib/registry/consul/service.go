package consul

import (
	"sync"
	"sync/atomic"

	"github.com/MrZhangjicheng/kitdemo/registry"
)

// 该结构体是将从数据中心的服务列表做本地缓存
type serviceSet struct {
	serviceName string
	services    *atomic.Value
	lock        sync.RWMutex
	watcher     map[*watcher]struct{}
}

func (s *serviceSet) broadcast(ss []*registry.ServiceIntance) {
	s.services.Store(ss)
	s.lock.RLock()
	defer s.lock.RUnlock()

	for k := range s.watcher {
		select {
		case k.event <- struct{}{}:
		default:
		}
	}
}
