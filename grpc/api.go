package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"antibrut/proto/antibrut/v1"
)

func (s *Server) Try(ctx context.Context, request *antibrut.TryRequest) (*antibrut.TryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Reset(ctx context.Context, request *antibrut.ResetRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) AddToWhiteList(ctx context.Context, request *antibrut.AddToWhiteListRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) DeleteFromWhiteList(ctx context.Context, request *antibrut.DeleteFromWhiteListRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) AddToBlackList(ctx context.Context, request *antibrut.AddToBlackListRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) DeleteFromBlackList(ctx context.Context, request *antibrut.DeleteFromBlackListRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
