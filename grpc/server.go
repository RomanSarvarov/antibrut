package grpc

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/romsar/antibrut"
	proto "github.com/romsar/antibrut/proto/antibrut/v1"
)

type Server struct {
	server  *grpc.Server
	service service

	proto.UnimplementedAntiBrutServiceServer
}

var _ proto.AntiBrutServiceServer = (*Server)(nil)

func NewServer(s service) *Server {
	return &Server{
		service: s,
	}
}

type service interface {
	Check(ctx context.Context, login antibrut.Login, pass antibrut.Password, ip antibrut.IP) error
	Reset(ctx context.Context, login antibrut.Login, ip antibrut.IP) error
	AddIPToWhiteList(ctx context.Context, subnet antibrut.Subnet) error
	DeleteIPFromWhiteList(ctx context.Context, subnet antibrut.Subnet) error
	AddIPToBlackList(ctx context.Context, subnet antibrut.Subnet) error
	DeleteIPFromBlackList(ctx context.Context, subnet antibrut.Subnet) error
}

func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "start grpc server error")
	}

	s.server = grpc.NewServer()
	proto.RegisterAntiBrutServiceServer(s.server, s)

	return s.server.Serve(lis)
}

func (s *Server) Close() error {
	if s.server != nil {
		s.server.GracefulStop()
	}
	return nil
}
