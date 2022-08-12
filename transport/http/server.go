package http

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"

	"github.com/MrZhangjicheng/kitdemo/internal/endpoint"
	"github.com/MrZhangjicheng/kitdemo/internal/host"
	"github.com/MrZhangjicheng/kitdemo/transport"
)

var (
	_ transport.Server = (*Server)(nil)
)

type Server struct {
	//  启动http服务的基本条件
	*http.Server
	lis      net.Listener
	endpoint *url.URL
	network  string
	address  string
	err      error
	// 是否加密
	tlsConf *tls.Config
}

func NewServer() *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
	}
	srv.Server = &http.Server{}
	return srv

}

func (srv *Server) Start(ctx context.Context) error {
	if err := srv.listenAndEndpoint(); err != nil {
		return nil
	}
	err := srv.Serve(srv.lis)
	if errors.Is(err, http.ErrServerClosed) {
		return srv.err
	}

	return nil

}

func (srv *Server) Stop(ctx context.Context) error {
	return srv.Shutdown(ctx)
}

func (srv *Server) listenAndEndpoint() error {
	if srv.lis == nil {
		lis, err := net.Listen(srv.network, srv.address)
		if err != nil {
			srv.err = err
			return err
		}
		srv.lis = lis
	}
	if srv.endpoint == nil {
		addr, err := host.Extract(srv.address, srv.lis)
		if err != nil {
			srv.err = err
			return err
		}
		srv.endpoint = endpoint.NewEndpoint(endpoint.Scheme("http", srv.tlsConf != nil), addr)
	}
	return nil
}
