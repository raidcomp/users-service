package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/raidcomp/users-service/clients"
	"github.com/raidcomp/users-service/daos"
	pb "github.com/raidcomp/users-service/proto"
	"github.com/raidcomp/users-service/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamoDBClient := clients.NewDynamoDBClient(cfg)
	usersDAO := daos.NewUsersDAO(dynamoDBClient)

	usersServer := server.NewUsersServer(usersDAO)
	grpcServer := grpc.NewServer()

	pb.RegisterUsersServer(grpcServer, usersServer)

	port := 5785
	address := fmt.Sprintf("localhost:%d", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening on %s", address)
	grpcServer.Serve(lis)
}
