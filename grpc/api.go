package grpc

import (
	"context"
	"errors"

	"github.com/romsar/antibrut"
	proto "github.com/romsar/antibrut/proto/antibrut/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Check проверяет "хороший" ли запрос, или его следует отклонить.
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

// Reset удаляет бакеты из хранилища.
func (s *Server) Reset(ctx context.Context, req *proto.ResetRequest) (*emptypb.Empty, error) {
	if req.GetLogin() == "" && req.GetIp() == "" {
		return nil, status.Error(codes.InvalidArgument, "no data to reset")
	}

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

// AddIPToWhiteList добавляет IP адрес в белый список.
func (s *Server) AddIPToWhiteList(ctx context.Context, req *proto.AddIPToWhiteListRequest) (*emptypb.Empty, error) {
	if req.GetSubnet() == "" {
		return nil, status.Error(codes.InvalidArgument, "no subnet passed")
	}

	err := s.service.AddIPToWhiteList(
		ctx,
		antibrut.Subnet(req.GetSubnet()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// DeleteIPFromWhiteList удаляет IP адрес из белого списка.
func (s *Server) DeleteIPFromWhiteList(
	ctx context.Context,
	req *proto.DeleteIPFromWhiteListRequest,
) (*emptypb.Empty, error) {
	if req.GetSubnet() == "" {
		return nil, status.Error(codes.InvalidArgument, "no subnet passed")
	}

	err := s.service.DeleteIPFromWhiteList(
		ctx,
		antibrut.Subnet(req.GetSubnet()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// AddIPToBlackList добавляет IP адрес в чёрный список.
func (s *Server) AddIPToBlackList(ctx context.Context, req *proto.AddIPToBlackListRequest) (*emptypb.Empty, error) {
	if req.GetSubnet() == "" {
		return nil, status.Error(codes.InvalidArgument, "no subnet passed")
	}

	err := s.service.AddIPToBlackList(
		ctx,
		antibrut.Subnet(req.GetSubnet()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// DeleteIPFromBlackList удаляет IP адрес из чёрного списка.
func (s *Server) DeleteIPFromBlackList(
	ctx context.Context,
	req *proto.DeleteIPFromBlackListRequest,
) (*emptypb.Empty, error) {
	if req.GetSubnet() == "" {
		return nil, status.Error(codes.InvalidArgument, "no subnet passed")
	}

	err := s.service.DeleteIPFromBlackList(
		ctx,
		antibrut.Subnet(req.GetSubnet()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
