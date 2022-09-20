package http

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/MrZhangjicheng/kitdemo/internal/endpoint"
	"github.com/MrZhangjicheng/kitdemo/log"
	"github.com/MrZhangjicheng/kitdemo/registry"
	"github.com/MrZhangjicheng/kitdemo/selector"
)

type Target struct {
	Scheme    string
	Authority string
	Endpoint  string
}

func parseTarget(endpoint string, insecure bool) (*Target, error) {
	if !strings.Contains(endpoint, "://") {
		if insecure {
			endpoint = "http://" + endpoint
		} else {
			endpoint = "https://" + endpoint
		}
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	target := &Target{Scheme: u.Scheme, Authority: u.Host}
	if len(u.Path) > 1 {
		target.Endpoint = u.Path[1:]
	}
	return target, nil
}

// 负责监听服务端列表的更改
type resolver struct {
	rebalancer selector.Rebalancer
	target     *Target
	watcher    registry.Watcher
	insecure   bool
}

func newResolver(ctx context.Context, discovery registry.Discover, target *Target, rebalancer selector.Rebalancer, block, insecure bool) (*resolver, error) {
	watcher, err := discovery.Watch(ctx, target.Endpoint)
	if err != nil {
		return nil, err
	}
	r := &resolver{
		target:     target,
		watcher:    watcher,
		rebalancer: rebalancer,
		insecure:   insecure,
	}
	if block {
		done := make(chan error, 1)
		go func() {
			for {
				services, err := watcher.Next()
				if err != nil {
					done <- err
					return
				}
				if r.update(services) {
					done <- nil
					return
				}
			}
		}()

		select {
		case err := <-done:
			if err != nil {
				stopErr := watcher.Stop()
				if stopErr != nil {
					log.Errorf("failed to http client watch stop: %v, error: %+v", target, stopErr)
				}
				return nil, err
			}
		case <-ctx.Done():
			log.Errorf("http client watch service %v reaching context deadline!", target)
			stopErr := watcher.Stop()
			if stopErr != nil {
				log.Errorf("failed to http client watch stop: %v,error :%+v", target, stopErr)
			}
			return nil, ctx.Err()
		}

	}

	// 非阻塞
	go func() {
		for {
			services, err := watcher.Next()
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				log.Errorf("http client watch service %v got unexpected error := %v", target, err)
				time.Sleep(time.Second)
				continue
			}
			r.update(services)
		}
	}()

	return r, nil

}

func (r *resolver) update(services []*registry.ServiceIntance) bool {
	nodes := make([]selector.Node, 0)
	for _, ins := range services {
		ept, err := endpoint.ParseEndpoint(ins.Endpoints, endpoint.Scheme("http", !r.insecure))
		if err != nil {
			log.Errorf("Failed to parse (%v) discovery endpoint: %v error %v", r.target, ins.Endpoints, err)
			continue
		}
		if ept == "" {
			continue
		}
		nodes = append(nodes, selector.NewNode("http", ept, ins))

	}
	if len(nodes) == 0 {
		log.Warnf("[http resolver]Zero endpoint found ,refused to write,set: %s ins: %v", r.target.Endpoint, nodes)
		return false
	}
	//
	r.rebalancer.Apply(nodes)

	return true

}
