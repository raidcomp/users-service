package server

import (
	"context"
	"github.com/raidcomp/users-service/daos"
	pb "github.com/raidcomp/users-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type usersServerImpl struct {
	pb.UnimplementedUsersServer

	UsersDAO daos.UsersDAO
}

func NewUsersServer(usersDAO daos.UsersDAO) pb.UsersServer {
	return usersServerImpl{
		UsersDAO: usersDAO,
	}
}

func (u usersServerImpl) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Check that login is not in use first
	userByLogin, err := u.UsersDAO.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error checking if login already exists")
	}

	if userByLogin != nil {
		return nil, status.Errorf(codes.AlreadyExists, "a user with login %s already exists", req.Login)
	}

	newUser, err := u.UsersDAO.CreateUser(ctx, req.Login, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating user")
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:        newUser.UserID,
			Login:     newUser.Login,
			Email:     newUser.Email,
			CreatedAt: timestamppb.New(newUser.CreatedAt),
			UpdatedAt: timestamppb.New(newUser.UpdatedAt),
		},
	}, nil
}

func (u usersServerImpl) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var user *daos.User
	var err error
	if req.Id != "" {
		user, err = u.UsersDAO.GetUserByID(ctx, req.Id)
	} else if req.Login != "" {
		user, err = u.UsersDAO.GetUserByLogin(ctx, req.Login)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting user")
	}

	if user == nil {
		return &pb.GetUserResponse{}, nil
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:        user.UserID,
			Login:     user.Login,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}
