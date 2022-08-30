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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
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
	// 拦截器 (grpc代码中只支持一个)
	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor
	// grpc 服务启动参数设置
	grpcOpts []grpc.ServerOption
}

func NewServer() *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
	}
	unaryInts := []grpc.UnaryServerInterceptor{}
	streamInts := []grpc.StreamServerInterceptor{}

	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}
	if len(srv.streamInts) > 0 {
		streamInts = append(streamInts, srv.streamInts...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
		grpc.ChainStreamInterceptor(streamInts...),
	}
	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	reflection.Register(srv.Server)
	return srv
}

func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
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
