package grpc

import (
	"context"
	"go-marketplace/internal/core/domain"
	users_service "go-marketplace/internal/features/users/service"
	users_grpc "go-marketplace/proto/gen"
)

type userServiceServer struct {
	users_grpc.UnimplementedUserServiceServer
	users *users_service.UserService
}

func NewUserServiceServer(users *users_service.UserService) *userServiceServer {
	return &userServiceServer{
		users: users,
	}
}

func (s *userServiceServer) CreateUser(ctx context.Context, req *users_grpc.CreateUserRequest) (*users_grpc.CreateUserResponse, error) {
	user, err := s.users.CreateUser(ctx, domain.User{
		Name:  req.GetName(),
		Phone: req.Phone,
		Role:  domain.UserRole(req.GetRole()),
	})
	if err != nil {
		return nil, err
	}

	return &users_grpc.CreateUserResponse{
		Id:   int64(user.ID),
		Role: string(user.Role),
	}, nil
}

func (s *userServiceServer) GetUserRole(ctx context.Context, req *users_grpc.GetUserRoleRequest) (*users_grpc.GetUserRoleResponse, error) {
	user, err := s.users.GetUser(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &users_grpc.GetUserRoleResponse{
		Role: string(user.Role),
	}, nil
}
