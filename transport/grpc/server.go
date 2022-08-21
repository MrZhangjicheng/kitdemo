package grpc

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"

	"github.com/MrZhangjicheng/kitdemo/internal/endpoint"
	"github.com/MrZhangjicheng/kitdemo/internal/host"
	"github.com/MrZhangjicheng/kitdemo/log"
	"github.com/MrZhangjicheng/kitdemo/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	_ transport.Server = (*Server)(nil)
)

//  grpc 服务对应的结构

type Server struct {
	// grpc 基本设置，可以成功启动
	*grpc.Server
	lis      net.Listener
	network  string
	address  string
	endpoint *url.URL
	err      error
	// 是否加密
	tlsConf *tls.Config
}

func NewServer() *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
	}
	grpcOpts := []grpc.ServerOption{}
	srv.Server = grpc.NewServer(grpcOpts...)
	reflection.Register(srv.Server)
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return s.err
	}
	log.Infof("[gRPC] server listening on: %s", s.lis.Addr().String())

	return s.Serve(s.lis)
}

func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	log.Info("[gRPC] server stopping")
	return nil
}

func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return err
		}
		s.lis = lis
	}
	if s.endpoint == nil {
		// 该函数主要是将serve服务运行的位置进行暴露，本地情况下就是本地地址加上端口
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			s.err = err
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("grpc", s.tlsConf != nil), addr)

	}
	return s.err
}
