package consul

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MrZhangjicheng/kitdemo/log"
	"github.com/MrZhangjicheng/kitdemo/registry"
	"github.com/hashicorp/consul/api"
)

type Client struct {
	cli    *api.Client
	ctx    context.Context
	cancel context.CancelFunc

	// 将从consul拿到的数据进行转化
	resolver ServiceResolver
	// 健康检测的间隔 和 通过心跳来检查服务是否可用 支持多种类型
	healthcheckInterval int
	// 是否开启心跳检测
	heartbeat bool
	// 多长时间取消
	deregisterCriticalServiceAfter int
}

type ServiceResolver func(ctx context.Context, entries []*api.ServiceEntry) []*registry.ServiceIntance

func defaultResolver(_ context.Context, entries []*api.ServiceEntry) []*registry.ServiceIntance {
	services := make([]*registry.ServiceIntance, 0, len(entries))

	for _, entry := range entries {
		// 版本号
		var version string
		for _, tag := range entry.Service.Tags {
			ss := strings.SplitN(tag, "=", 2)
			if len(ss) == 2 && ss[0] == "version" {
				version = ss[1]
			}
		}
		endpoints := make([]string, 0)
		for scheme, addr := range entry.Service.TaggedAddresses {
			if scheme == "lan_ipv4" || scheme == "wan_ipv4" || scheme == "lan_ipv6" || scheme == "wan_ipv6" {
				continue
			}
			endpoints = append(endpoints, addr.Address)
		}
		if len(endpoints) == 0 && entry.Service.Address != "" && entry.Service.Port != 0 {
			endpoints = append(endpoints, fmt.Sprintf("http://%s:%d", entry.Service.Address, entry.Service.Port))
		}

		services = append(services, &registry.ServiceIntance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Endpoints: endpoints,
			Version:   version,
		})
	}

	return services
}

func NewClient(cli *api.Client) *Client {
	c := &Client{
		cli:                            cli,
		resolver:                       defaultResolver,
		healthcheckInterval:            10,
		deregisterCriticalServiceAfter: 600,
		heartbeat:                      true,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c
}

// 注册
func (c *Client) Register(_ context.Context, svc *registry.ServiceIntance, enableHealthCheck bool) error {
	// 构建consul中的服务真实地址
	address := make(map[string]api.ServiceAddress)
	// 用来健康检查
	checkAddresses := make([]string, 0, len(svc.Endpoints))
	for _, endpoint := range svc.Endpoints {
		raw, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		addr := raw.Hostname()
		port, _ := strconv.ParseUint(raw.Port(), 10, 16)
		checkAddresses = append(checkAddresses, net.JoinHostPort(addr, strconv.FormatUint(port, 10)))
		address[raw.Scheme] = api.ServiceAddress{Address: endpoint, Port: int(port)}

	}
	asr := &api.AgentServiceRegistration{
		ID:              svc.ID,
		Name:            svc.Name,
		TaggedAddresses: address,
	}
	if len(checkAddresses) > 0 {
		host, portRaw, _ := net.SplitHostPort(checkAddresses[0])
		port, _ := strconv.ParseInt(portRaw, 10, 32)
		asr.Address = host
		asr.Port = int(port)
	}
	if enableHealthCheck {
		for _, address := range checkAddresses {
			asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
				TCP:                            address,
				Interval:                       fmt.Sprintf("%ds", c.healthcheckInterval),
				DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", c.deregisterCriticalServiceAfter),
				Timeout:                        "5s",
			})
		}

	}
	if c.heartbeat {
		asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
			CheckID:                        "service:" + svc.ID,
			TTL:                            fmt.Sprintf("%ds", c.healthcheckInterval*2),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", c.deregisterCriticalServiceAfter),
		})

	}
	err := c.cli.Agent().ServiceRegister(asr)
	if err != nil {
		return err

	}
	if c.heartbeat {
		go func() {
			time.Sleep(time.Second)
			err := c.cli.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass")
			if err != nil {
				log.Errorf("[Consul]update ttl heartbeat to consul failed!err:=%v", err)
			}
			ticker := time.NewTicker(time.Second * time.Duration(c.healthcheckInterval))
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					err := c.cli.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass")
					if err != nil {
						log.Errorf("[Consul]update ttl heartbeat to consul failed!err:=%v", err)
					}
				case <-c.ctx.Done():
					return
				}
			}
		}()
	}
	return nil
}

// 销毁
func (c *Client) Deregister(_ context.Context, serviceID string) error {
	c.cancel()
	return c.cli.Agent().ServiceDeregister(serviceID)
}

// 这个复杂查找采用了类似分页的方法 先不用
func (c *Client) Service(ctx context.Context, service string, index uint64, passingOnly bool) ([]*registry.ServiceIntance, uint64, error) {
	opts := &api.QueryOptions{
		WaitIndex: index,
		WaitTime:  time.Second * 55,
	}
	opts = opts.WithContext(ctx)
	entries, meta, err := c.cli.Health().Service(service, "", passingOnly, opts)
	if err != nil {
		return nil, 0, err
	}

	return c.resolver(ctx, entries), meta.LastIndex, nil
}
