package main

import (
	pb "br-prep/proto"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func getDBURL() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	return url
}

func openDB(ctx context.Context, dbURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err = db.PingContext(pingCtx)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}
	return db, nil
}

func main() {
	ctx := context.Background()
	dbURL := getDBURL()

	db, err := openDB(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error %v occurred", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()

	// Infrastructure (repo)
	// userRepo := NewInMemoryUserRepository()
	// _ = userRepo
	repo := NewPostgresUserRepository(db)

	// Services (API layer)
	healthSvc := &healthServiceServer{}
	userSvc := NewUserServiceServer(repo)

	pb.RegisterHealthServiceServer(grpcServer, healthSvc)
	pb.RegisterUserServiceServer(grpcServer, userSvc)

	log.Println("gRPC(TCP) Server is listening on port 50051...")
	log.Fatal(grpcServer.Serve(lis))
}
