package main

import (
	pb "br-prep/proto"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	repo UserRepository
}

func NewUserServiceServer(repo UserRepository) *userServiceServer {
	return &userServiceServer{repo: repo}
}

func (s *userServiceServer) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	id := req.GetId()
	if id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	u, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.UserResponse{Name: u.Name}, nil
}

func (s *userServiceServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersResponse, error) {
	ids := req.GetIds()
	if len(ids) == 0 {
		return &pb.BatchGetUsersResponse{Users: []*pb.UserResponse{}}, nil
	}

	users, err := s.repo.BatchGetUsers(ctx, ids)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	out := make([]*pb.UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, &pb.UserResponse{Name: u.Name})
	}

	return &pb.BatchGetUsersResponse{Users: out}, nil
}
