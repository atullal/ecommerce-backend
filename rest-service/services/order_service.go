package services

import (
    "context"
    pb "github.com/atullal/ecommerce-backend-protobuf/order" // Replace with the correct import path
)

type OrderService struct {
    GrpcClient pb.OrderServiceClient
}

func NewOrderService(client pb.OrderServiceClient) *OrderService {
    return &OrderService{
        GrpcClient: client,
    }
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
    return s.GrpcClient.CreateOrder(ctx, req)
}

func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
    return s.GrpcClient.GetOrder(ctx, req)
}

func (s *OrderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {
    return s.GrpcClient.UpdateOrder(ctx, req)
}

func (s *OrderService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
    return s.GrpcClient.ListOrders(ctx, req)
}
