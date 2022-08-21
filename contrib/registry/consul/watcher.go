package consul

import (
	"context"

	"github.com/MrZhangjicheng/kitdemo/registry"
)

type watcher struct {
	event chan struct{}

	set *serviceSet

	ctx context.Context

	cancel context.CancelFunc
}

func (w *watcher) Next() (service []*registry.ServiceIntance, err error) {
	select {
	case <-w.ctx.Done():
		err = w.ctx.Err()
	case <-w.event:
	}

	ss, ok := w.set.services.Load().([]*registry.ServiceIntance)
	if ok {
		service = append(service, ss...)
	}
	return
}

func (w *watcher) Stop() error {
	w.cancel()
	w.set.lock.Lock()
	defer w.set.lock.Unlock()
	delete(w.set.watcher, w)
	return nil
}
