package main

import (
	"context"
	"fmt"
	pb "github.com/Dr-Evans/raidcomp/users-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

type usersServer struct {
	pb.UnimplementedUsersServer
}

func (u usersServer) GetUser(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return &pb.GetUserResponse{
		User: &pb.User{
			Id:        "test-id",
			Login:     "test-login",
			Email:     "test-email",
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		},
	}, nil
}

func main() {
	port := 5785
	address := fmt.Sprintf("localhost:%d", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUsersServer(grpcServer, &usersServer{})

	log.Printf("Listening on %s", address)
	grpcServer.Serve(lis)
}
