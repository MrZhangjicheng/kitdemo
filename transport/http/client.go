package http

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/MrZhangjicheng/kitdemo/registry"
)

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
}

type Client struct {
	opts     clientOptions
	cc       *http.Client
	insecure bool
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
	if target.Scheme == "discovery" {
		watcher, err := options.discovery.Watch(context.Background(), target.Endpoint)
		if err != nil {
			panic("服务发现错误")
		}
		srvs, _ := watcher.Next()
		// 同时开启http服务和grpc服务选择对服务
		srv := srvs[0]
		for _, k := range srv.Endpoints {
			if strings.Contains(k, "http") {
				options.endpoint = k[7:]
			}
		}
	}

	return &Client{
		opts:     options,
		insecure: insecure,
		cc:       &http.Client{},
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
