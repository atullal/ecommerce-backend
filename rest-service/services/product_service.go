package services

import (
    "context"
    pb "github.com/atullal/ecommerce-backend-protobuf/product" // Replace with the correct import path
)

type ProductService struct {
    GrpcClient pb.ProductServiceClient
}

func NewProductService(client pb.ProductServiceClient) *ProductService {
    return &ProductService{
        GrpcClient: client,
    }
}

func (s *ProductService) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.ProductResponse, error) {
    return s.GrpcClient.AddProduct(ctx, req)
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
    return s.GrpcClient.GetProduct(ctx, req)
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
    return s.GrpcClient.UpdateProduct(ctx, req)
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
    return s.GrpcClient.DeleteProduct(ctx, req)
}

func (s *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
    return s.GrpcClient.ListProducts(ctx, req)
}

func (s *ProductService) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryRequest) (*pb.InventoryResponse, error) {
    return s.GrpcClient.UpdateInventory(ctx, req)
}

func (s *ProductService) GetInventory(ctx context.Context, req *pb.GetInventoryRequest) (*pb.InventoryResponse, error) {
    return s.GrpcClient.GetInventory(ctx, req)
}
