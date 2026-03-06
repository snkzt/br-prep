package main

import (
	"context"

	pb "br-prep/proto"
)

type healthServiceServer struct {
	pb.UnimplementedHealthServiceServer
}

func (hs *healthServiceServer) Check(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Message: "OK"}, nil
}
