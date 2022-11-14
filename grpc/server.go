package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/romsar/antibrut"
	proto "github.com/romsar/antibrut/proto/antibrut/v1"
	"google.golang.org/grpc"
)

// Server предоставляет структура GRPC сервера.
type Server struct {
	// server содержит сущность GRPC сервера.
	server *grpc.Server

	// service содержит методы для работы с бизнес-логикой.
	service service

	// logger механизм логирования.
	logger logger

	proto.UnimplementedAntiBrutServiceServer
}

var _ proto.AntiBrutServiceServer = (*Server)(nil)

// Option возвращает функцию, модифицирующую Server.
type Option func(s *Server)

// WithLogger возвращает функцию,
// устанавливающую механизм логирования.
func WithLogger(l logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}

// NewServer создает Server.
func NewServer(s service, opts ...Option) *Server {
	srv := &Server{
		service: s,
	}

	for _, opt := range opts {
		opt(srv)
	}

	if srv.logger == nil {
		srv.logger = log.New(io.Discard, "", 0)
	}

	return srv
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

// logger это контракт для механизма логирования.
type logger interface {
	// Printf сохраняет отформатированное сообщение в лог.
	Printf(format string, v ...any)
}

// Start запускает Server и слушает входящие соединения.
func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("start grpc server error: %w", err)
	}

	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(s.LoggingInterceptor),
	)

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

// LoggingInterceptor перехватчик запросов, добавляющий логирование.
func (s *Server) LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	rsp, err := handler(ctx, req)

	s.logger.Printf(
		"grpc: method=%s\tduration=%s\terror=%v\treq=%v\trsp%v\n",
		info.FullMethod,
		time.Since(start),
		err,
		req,
		rsp,
	)

	return rsp, err
}
