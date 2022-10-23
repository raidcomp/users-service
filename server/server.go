package server

import (
	"context"
	"errors"
	"github.com/raidcomp/users-service/auth"
	"github.com/raidcomp/users-service/daos"
	pb "github.com/raidcomp/users-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"unicode"
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

func validateLogin(login string) error {
	if login == "" {
		return errors.New("login must not be empty")
	}

	var (
		hasNotAllowed                  = false
		hasUnderscoreAtStart           = false
		hasUnderscoreAtEnd             = false
		hasPeriodAtStart               = false
		hasPeriodAtEnd                 = false
		hasUnderscoreAndPeriodTogether = false
	)

	for i, char := range login {
		if !unicode.IsLetter(char) || !unicode.IsNumber(char) || !(char == '_') || !(char == '.') {
			hasNotAllowed = true
			break
		}

		if i == 0 {
			if char == '_' {
				hasUnderscoreAtStart = true
				break
			} else if char == '.' {
				hasPeriodAtStart = true
				break
			}
		}

		if i == len(login)-1 {
			if char == '_' {
				hasUnderscoreAtEnd = true
				break
			} else if char == '.' {
				hasPeriodAtEnd = true
				break
			}
		}

		if i < len(login)-1 {
			nextChar := login[i+1]
			if (char == '_' || char == '.') && (nextChar == '_' || nextChar == '.') {
				hasUnderscoreAndPeriodTogether = true
				break
			}
		}
	}

	if !hasNotAllowed {
		return errors.New("login must only contain letters or numbers")
	} else if !hasUnderscoreAtStart {
		return errors.New("login must not start with an underscore")
	} else if !hasUnderscoreAtEnd {
		return errors.New("login must not end with an underscore")
	} else if !hasPeriodAtStart {
		return errors.New("login must not start with a period")
	} else if !hasPeriodAtEnd {
		return errors.New("login must not end with a period")
	} else if !hasUnderscoreAndPeriodTogether {
		return errors.New("login must not have two consecutive underscores or periods")
	} else {
		return nil
	}
}

func validatePassword(password string) error {
	if password == "" {
		return errors.New("password must not be empty")
	}

	var (
		hasUpper      = false
		hasLower      = false
		hasNumber     = false
		hasSpecial    = false
		hasNotAllowed = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}

		if !unicode.IsLetter(char) || !unicode.IsNumber(char) || !unicode.IsPunct(char) || !unicode.IsSymbol(char) {
			hasNotAllowed = true
			break
		}
	}

	if !hasNotAllowed {
		return errors.New("password must only contain letter, numbers, or special characters")
	} else if !hasUpper {
		return errors.New("password must contain uppercase character")
	} else if !hasLower {
		return errors.New("password must contain lowercase character")
	} else if !hasNumber {
		return errors.New("password must contain a number")
	} else if !hasSpecial {
		return errors.New("password must contain a special character")
	} else {
		return nil
	}
}

func (u usersServerImpl) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	err = validatePassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "password invalid")
	}

	err = validateLogin(req.Login)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "login invalid")
	}

	// Check that login is not in use first
	userByLogin, err := u.UsersDAO.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error checking if login already exists")
	}

	if userByLogin != nil {
		return nil, status.Errorf(codes.AlreadyExists, "a user with login %s already exists", req.Login)
	}

	newUser, err := u.UsersDAO.CreateUser(ctx, req.Login, req.Email, req.Password)
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
	err := req.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var user *daos.User
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

func (u usersServerImpl) CheckUserPassword(ctx context.Context, req *pb.CheckUserPasswordRequest) (*pb.CheckUserPasswordResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var user *daos.User
	if req.Id != "" {
		user, err = u.UsersDAO.GetUserByID(ctx, req.Id)
	} else if req.Login != "" {
		user, err = u.UsersDAO.GetUserByLogin(ctx, req.Login)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting user")
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user for userID %s not found", user.UserID)
	}

	if !auth.CheckPasswordHash(user.HashedPassword, req.Password) {
		return nil, status.Errorf(codes.InvalidArgument, "password does not match userID %s password", user.UserID)
	}

	return &pb.CheckUserPasswordResponse{}, nil
}
