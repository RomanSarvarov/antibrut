package grpc

import (
	"antibrut/proto/antibrut/v1"
)

type Server struct {
	antibrut.UnimplementedAntiBrutServiceServer
}

var _ antibrut.AntiBrutServiceServer = (*Server)(nil)

func NewServer() *Server {
	return &Server{}
}
