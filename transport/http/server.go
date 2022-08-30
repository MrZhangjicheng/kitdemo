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
	"github.com/MrZhangjicheng/kitdemo/log"
	"github.com/MrZhangjicheng/kitdemo/transport"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
	_ http.Handler         = (*Server)(nil)
)

type ServerOption func(*Server)

func WithHandler(handler http.Handler) ServerOption {
	return func(s *Server) {
		s.handler = handler
	}
}

// http 服务应该具备的能力
// 路由注册
// 中间件
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
	// 路由能力 支持动态路由以及路由组的概念
	// router      *mux.Router
	// strictSlash bool // 路由规则是否严格
	handler http.Handler
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		// strictSlash: true,
	}
	// srv.router = mux.NewRouter().StrictSlash(srv.strictSlash)
	// srv.router.NotFoundHandler = http.DefaultServeMux
	// srv.router.MethodNotAllowedHandler = http.DefaultServeMux
	for _, o := range opts {
		o(srv)
	}
	if srv.handler == nil {
		panic("路由未注册")
	}
	srv.Server = &http.Server{
		Handler: srv.handler,
	}
	return srv

}

// 路由相关
// func (s *Server) Route(prefix string) *Router {
// 	return newRouter(prefix)
// }

func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}

	return s.endpoint, nil
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// httpServer 中的handler 接口
	s.Handler.ServeHTTP(res, req)
}

func (srv *Server) Start(ctx context.Context) error {
	if err := srv.listenAndEndpoint(); err != nil {
		return nil
	}
	log.Infof("[http] server listening on: %s", srv.lis.Addr().String())
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
