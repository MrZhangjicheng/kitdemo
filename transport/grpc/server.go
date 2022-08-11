package grpc

import (
	"context"

	"github.com/MrZhangjicheng/kitdemo/transport"
	"google.golang.org/grpc"
)

var (
	_ transport.Server = (*Server)(nil)
)

type Server struct {
	*grpc.Server
	address string
}

func (s *Server) Start(ctx context.Context) error {
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}
