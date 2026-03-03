package main

import (
	"context"
	"log"
	"time"

	pb "br-prep/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	clientHealth := pb.NewHealthServiceClient(conn)
	clientUser := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	respHealth, err := clientHealth.Check(ctx, &pb.HealthRequest{})
	if err != nil {
		log.Fatal(err)
	}

	respUser, err := clientUser.GetUser(ctx, &pb.UserRequest{Id: 1})
	if err != nil {
		log.Fatal(err)
	}

	respUsers, err := clientUser.BatchGetUsers(ctx, &pb.BatchGetUsersRequest{Ids: []int64{1, 2, 999}})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(respHealth.Message)
	log.Println(respUser.Name)
	log.Println(respUsers.Users)
}
