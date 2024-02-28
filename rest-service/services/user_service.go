package services

import (
    "context"
    pb "github.com/atullal/ecommerce-backend/user-service/gen" // Replace with the correct import path
)

type UserService struct {
    GrpcClient pb.UserServiceClient
}

func NewUserService(client pb.UserServiceClient) *UserService {
    return &UserService{
        GrpcClient: client,
    }
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
    return s.GrpcClient.CreateUser(ctx, req)
}

func (s *UserService) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.UserResponse, error) {
    return s.GrpcClient.AuthenticateUser(ctx, req)
}

// Additional business logic functions can be added here...
