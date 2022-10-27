package grpc

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/romsar/antibrut"
	proto "github.com/romsar/antibrut/proto/antibrut/v1"
)

// Server предоставляет структура GRPC сервера.
type Server struct {
	server  *grpc.Server
	service service

	proto.UnimplementedAntiBrutServiceServer
}

var _ proto.AntiBrutServiceServer = (*Server)(nil)

// NewServer создает Server.
func NewServer(s service) *Server {
	return &Server{
		service: s,
	}
}

// service декларирует методы, необходимые Server.
type service interface {
	// Check проверяет "хороший" ли запрос, или его следует отклонить.
	Check(ctx context.Context, login antibrut.Login, pass antibrut.Password, ip antibrut.IP) error

	// Reset удаляет бакеты из хранилища.
	Reset(ctx context.Context, login antibrut.Login, ip antibrut.IP) error

	// AddIPToWhiteList добавляет IP адрес в белый список.
	AddIPToWhiteList(ctx context.Context, subnet antibrut.Subnet) error

	// DeleteIPFromWhiteList удаляет IP адрес из белого списка.
	DeleteIPFromWhiteList(ctx context.Context, subnet antibrut.Subnet) error

	// AddIPToBlackList добавляет IP адрес в чёрный список.
	AddIPToBlackList(ctx context.Context, subnet antibrut.Subnet) error

	// DeleteIPFromBlackList удаляет IP адрес из чёрного списка.
	DeleteIPFromBlackList(ctx context.Context, subnet antibrut.Subnet) error
}

// Start запускает Server и слушает входящие соединения.
func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "start grpc server error")
	}

	s.server = grpc.NewServer()
	proto.RegisterAntiBrutServiceServer(s.server, s)

	return s.server.Serve(lis)
}

// Close останавливает Server.
func (s *Server) Close() error {
	if s.server != nil {
		s.server.GracefulStop()
	}
	return nil
}
