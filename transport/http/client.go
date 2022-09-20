package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/MrZhangjicheng/kitdemo/internal/host"

	"github.com/MrZhangjicheng/kitdemo/registry"
	"github.com/MrZhangjicheng/kitdemo/selector"
	"github.com/MrZhangjicheng/kitdemo/selector/random"
)

// 服务启动时初始化默认的全局selector的构建器
func init() {
	if selector.GlobalSelector() == nil {
		selector.SetGlobalSelector(random.NewBuilder())
	}
}

type ClientOption func(*clientOptions)

func WithEndpoint(endpoint string) ClientOption {
	return func(co *clientOptions) {
		co.endpoint = endpoint
	}
}

func Withdiscovery(d registry.Discover) ClientOption {
	return func(co *clientOptions) {
		co.discovery = d
	}
}

type clientOptions struct {
	ctx       context.Context
	tlsConf   *tls.Config
	endpoint  string
	discovery registry.Discover
	block     bool
}

type Client struct {
	opts     clientOptions
	target   *Target
	cc       *http.Client
	insecure bool
	r        *resolver
	selector selector.Selector
}

func Newclient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	options := clientOptions{
		ctx: ctx,
	}
	for _, o := range opts {
		o(&options)
	}
	insecure := options.tlsConf == nil
	target, err := parseTarget(options.endpoint, insecure)
	if err != nil {
		return nil, err
	}
	selector := selector.GlobalSelector().Build()
	var r *resolver
	if options.discovery != nil {
		if target.Scheme == "discovery" {
			if r, err = newResolver(ctx, options.discovery, target, selector, options.block, insecure); err != nil {
				return nil, fmt.Errorf("[http client] new resolver failed!err: %v", options.endpoint)
			} else if _, _, err := host.ExtractHostPort(options.endpoint); err != nil {
				return nil, fmt.Errorf("[http client] invalid endpoint format: %v", options.endpoint)
			}
		}
	}
	return &Client{
		opts:     options,
		insecure: insecure,
		r:        r,
		cc:       &http.Client{},
		selector: selector,
	}, nil

}

func (client *Client) Do(method, path string) (*http.Response, error) {
	url := "http"
	if !client.insecure {
		url += "s"
	}
	url += "://" + client.opts.endpoint + path
	req, err := http.NewRequest(method, url, nil)
	resp, err := client.cc.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
