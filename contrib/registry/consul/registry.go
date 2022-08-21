package consul

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MrZhangjicheng/kitdemo/registry"
	"github.com/hashicorp/consul/api"
)

var _ registry.Registar = &Registry{}
var _ registry.Discover = &Registry{}

// 注册销毁 需要考虑服务发现
type Registry struct {
	cli               *Client
	enableHealthCheck bool
	registry          map[string]*serviceSet
	lock              sync.RWMutex
}

func New(apiClient *api.Client) *Registry {
	r := &Registry{
		cli:               NewClient(apiClient),
		enableHealthCheck: true,
		registry:          make(map[string]*serviceSet),
	}

	return r

}

func (r *Registry) Register(ctx context.Context, svc *registry.ServiceIntance) error {
	return r.cli.Register(ctx, svc, r.enableHealthCheck)
}

func (r *Registry) Deregister(ctx context.Context, svc *registry.ServiceIntance) error {
	return r.cli.Deregister(ctx, svc.ID)
}

// 先从本地缓存中找
func (r *Registry) FindService(ctx context.Context, name string) ([]*registry.ServiceIntance, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	set := r.registry[name]

	getRemote := func() []*registry.ServiceIntance {
		services, _, err := r.cli.Service(ctx, name, 0, true)
		if err == nil && len(services) > 0 {
			return services
		}
		return nil
	}
	if set == nil {
		if s := getRemote(); len(s) > 0 {
			return s, nil
		}
		return nil, fmt.Errorf("service %s not resolved in registry", name)
	}
	ss, _ := set.services.Load().([]*registry.ServiceIntance)
	if ss == nil {
		if s := getRemote(); len(s) > 0 {
			return s, nil
		}
		return nil, fmt.Errorf("service %s not resolved in registry", name)
	}
	return ss, nil

}

func (r *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	set, ok := r.registry[name]

	if !ok {
		set = &serviceSet{
			watcher:     make(map[*watcher]struct{}),
			services:    &atomic.Value{},
			serviceName: name,
		}
		r.registry[name] = set
	}
	w := &watcher{
		event: make(chan struct{}, 1),
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.set = set
	set.lock.Lock()
	set.watcher[w] = struct{}{}

	set.lock.Unlock()

	ss, _ := set.services.Load().([]*registry.ServiceIntance)

	if len(ss) > 0 {
		w.event <- struct{}{}
	}
	if !ok {
		err := r.resolve(set)
		if err != nil {
			return nil, err
		}
	}

	return w, nil

}

func (r *Registry) resolve(ss *serviceSet) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	services, idx, err := r.cli.Service(ctx, ss.serviceName, 0, true)
	cancel()
	if err != nil {
		return err
	} else if len(services) > 0 {
		ss.broadcast(services)
	}
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
			tmpService, tmpIdx, err := r.cli.Service(ctx, ss.serviceName, idx, true)
			cancel()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			if len(tmpService) != 0 && tmpIdx != idx {
				services = tmpService
				ss.broadcast(services)
			}
			idx = tmpIdx
		}
	}()
	return nil

}
