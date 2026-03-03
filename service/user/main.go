package main

import (
	pb "br-prep/proto"
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type healthServiceServer struct {
	pb.UnimplementedHealthServiceServer
}

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	mu    sync.RWMutex
	users map[int64]*pb.UserResponse
}

func (hs *healthServiceServer) Check(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Message: "OK"}, nil
}

func newUserServer() *userServiceServer {
	return &userServiceServer{
		users: map[int64]*pb.UserResponse{
			1:  {Name: "Natalie"},
			2:  {Name: "Domingo"},
			11: {Name: "Lucy"},
			99: {Name: "Rasicov"},
		},
	}
}

func (u *userServiceServer) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	id := req.GetId()

	if id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	u.mu.RLock()
	user, ok := u.users[id]
	u.mu.RUnlock()
	if !ok {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return user, nil
}

func (u *userServiceServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersResponse, error) {
	ids := req.GetIds()

	var names []*pb.UserResponse
	for _, id := range ids {
		u.mu.RLock()
		user, ok := u.users[id]
		u.mu.RUnlock()
		if ok {
			names = append(names, user)
		}
	}

	return &pb.BatchGetUsersResponse{
		Users: names,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error %v occurred", err)
	}
	defer lis.Close()

	log.Println("TCP Server is listening on port 50051...")

	grpcServer := grpc.NewServer()
	pb.RegisterHealthServiceServer(grpcServer, &healthServiceServer{})
	pb.RegisterUserServiceServer(grpcServer, newUserServer())
	log.Fatal(grpcServer.Serve(lis))
}
