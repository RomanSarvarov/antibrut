package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/romsar/antibrut"
	proto "github.com/romsar/antibrut/proto/antibrut/v1"
)

func (s *Server) Check(ctx context.Context, req *proto.CheckRequest) (*proto.CheckResponse, error) {
	if req.GetLogin() == "" && req.GetPassword() == "" && req.GetIp() == "" {
		return nil, status.Error(codes.InvalidArgument, "no data to check")
	}

	err := s.service.Check(
		ctx,
		antibrut.Login(req.GetLogin()),
		antibrut.Password(req.GetPassword()),
		antibrut.IP(req.GetIp()),
	)

	if err != nil {
		if errors.Is(err, antibrut.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, antibrut.ErrMaxAttemptsExceeded) || errors.Is(err, antibrut.ErrIPInBlackList) {
			return &proto.CheckResponse{
				Ok: false,
			}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.CheckResponse{
		Ok: true,
	}, nil
}

func (s *Server) Reset(ctx context.Context, req *proto.ResetRequest) (*emptypb.Empty, error) {
	err := s.service.Reset(
		ctx,
		antibrut.Login(req.GetLogin()),
		antibrut.IP(req.GetIp()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) AddIPToWhiteList(ctx context.Context, req *proto.AddIPToWhiteListRequest) (*emptypb.Empty, error) {
	err := s.service.AddIPToWhiteList(
		ctx,
		antibrut.Subnet(req.GetSubnet()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteIPFromWhiteList(ctx context.Context, req *proto.DeleteIPFromWhiteListRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteIPFromWhiteList(
		ctx,
		antibrut.Subnet(req.GetSubnet()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) AddIPToBlackList(ctx context.Context, req *proto.AddIPToBlackListRequest) (*emptypb.Empty, error) {
	err := s.service.AddIPToBlackList(
		ctx,
		antibrut.Subnet(req.GetSubnet()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteIPFromBlackList(ctx context.Context, req *proto.DeleteIPFromBlackListRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteIPFromBlackList(
		ctx,
		antibrut.Subnet(req.GetSubnet()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
