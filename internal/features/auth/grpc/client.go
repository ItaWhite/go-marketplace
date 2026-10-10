package auth

import (
	"context"
	"fmt"
	"go-marketplace/internal/core/transport/grpc_errors"
	"go-marketplace/internal/features/auth/service"
	"go-marketplace/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userServiceClient struct {
	client users_grpc.UserServiceClient
}

func NewUserServiceClient(addr string) (*userServiceClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error creating grpc connection: %w", err)
	}

	return &userServiceClient{
		client: users_grpc.NewUserServiceClient(conn),
	}, nil
}

func (c *userServiceClient) Create(ctx context.Context, input service.CreateUserInput) (int, string, error) {
	req := &users_grpc.CreateUserRequest{
		Name:  input.Name,
		Phone: input.Phone,
		Role:  input.Role,
	}

	resp, err := c.client.CreateUser(ctx, req)
	if err != nil {
		return 0, "", grpc_errors.ToSentinelError(err)
	}

	return int(resp.GetId()), resp.GetRole(), nil
}

func (c *userServiceClient) GetRole(ctx context.Context, id int) (string, error) {
	req := &users_grpc.GetUserRoleRequest{
		Id: int64(id),
	}

	resp, err := c.client.GetUserRole(ctx, req)
	if err != nil {
		return "", grpc_errors.ToSentinelError(err)
	}

	return resp.GetRole(), nil
}
