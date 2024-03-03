package main

import (
	"os"
    "context"
    "log"
    "net"
    "errors"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    pb "github.com/atullal/ecommerce-backend-protobuf/order"
    productpb "github.com/atullal/ecommerce-backend-protobuf/product"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "order-service/models"
    "fmt"
    "time"
    "math"
)

type server struct {
    pb.OrderServiceServer
    db *gorm.DB
    ProductServiceClient productpb.ProductServiceClient
}
func connectWithBackoff(dsn string) (*gorm.DB, error) {
    var db *gorm.DB
    var err error

    maxAttempts := 5
    for i := 0; i < maxAttempts; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            sqlDB, err := db.DB()
            if err == nil {
                err = sqlDB.Ping() // Try to ping the database
                if err == nil {
                    return db, nil // Success
                }
            }
        }

        // Exponential backoff logic
        backoff := math.Pow(2, float64(i))
        time.Sleep(time.Duration(backoff) * time.Second)
    }

    return nil, err // Return the last error
}

func initDB() *gorm.DB {
    dsn := os.Getenv("DSN")
    db, err := connectWithBackoff(dsn)
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }

    // Migrate the schema
    if err := db.AutoMigrate(&models.Order{}, &models.OrderItem{}); err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }
    fmt.Println("Database connection successful")
    return db
}

func (s *server) connectToProductService() {

	// Set up a connection to the gRPC server.
    productServiceConnection, err := grpc.Dial("0.0.0.0:50052", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }
    fmt.Println("Connected to gRPC server")

    // Initialize the ProductService and ProductHandler
    s.ProductServiceClient = productpb.NewProductServiceClient(productServiceConnection)
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
    fmt.Println("Create order request", req)

    // Convert req items to order items and prepare inventory updates
    orderItems := make([]models.OrderItem, 0, len(req.Items))
    inventoriesReq := &productpb.UpdateMultipleInventoriesRequest{}

    for _, item := range req.Items {
        orderItem := models.OrderItem{
            ProductID: uint(item.ProductId),
            Quantity:  int(item.Quantity),
        }
        orderItems = append(orderItems, orderItem)

        inventoryReq := &productpb.UpdateInventoryRequest{
            ProductId: int64(item.ProductId),
            QuantityChange: -int32(item.Quantity), // Decrease the quantity
        }
        inventoriesReq.InventoryUpdates = append(inventoriesReq.InventoryUpdates, inventoryReq)
    }

    // Start a transaction
    tx := s.db.Begin()

    // Update inventory in product-service
    _, err := s.ProductServiceClient.UpdateMultipleInventories(ctx, inventoriesReq)
    if err != nil {
        tx.Rollback()
        return nil, status.Errorf(codes.Internal, "Error updating inventory: %v", err)
    }

    // Create the order in the database
    newOrder := models.Order{
        CustomerID: uint(req.CustomerId),
        Items:      orderItems,
        Status:     models.StatusPending,
        // Set other fields based on your request and models
    }

    if err := tx.Create(&newOrder).Error; err != nil {
        tx.Rollback()
        return nil, status.Errorf(codes.Internal, "Error creating order: %v", err)
    }

    // Commit the transaction
    tx.Commit()

    // Prepare and return the response
    response := &pb.OrderResponse{
        Order: &pb.Order{
            Id:         int64(newOrder.ID),
            CustomerId: int64(newOrder.CustomerID),
            // Populate other fields of the protobuf Order message
        },
    }

    return response, nil
}

func mapOrderStatusToProto(status models.OrderStatus) pb.OrderStatus {
    switch status {
    case models.StatusPending:
        return pb.OrderStatus_PENDING
    case models.StatusConfirmed:
        return pb.OrderStatus_CONFIRMED
    case models.StatusShipped:
        return pb.OrderStatus_SHIPPED
    case models.StatusDelivered:
        return pb.OrderStatus_DELIVERED
    case models.StatusCancelled:
        return pb.OrderStatus_CANCELLED
    default:
        // Handle default case or return an error
        return pb.OrderStatus_PENDING // Assuming you have an UNKNOWN status in your protobuf enum
    }
}

func (s *server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
    var order models.Order

    // Retrieve the order by ID from the database
    result := s.db.Preload("Items").First(&order, req.OrderId)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, status.Errorf(codes.NotFound, "Order with ID '%d' not found", req.OrderId)
        }
        return nil, status.Errorf(codes.Internal, "Error retrieving order: %v", result.Error)
    }

    // Convert the order and its items to the protobuf type
    orderItems := make([]*pb.OrderItem, len(order.Items))
    for i, item := range order.Items {
        orderItems[i] = &pb.OrderItem{
            ProductId: int64(item.ProductID),
            Quantity:  int32(item.Quantity),
        }
    }

    // Prepare and return the response
    response := &pb.OrderResponse{
        Order: &pb.Order{
            Id:         int64(order.ID),
            CustomerId: int64(order.CustomerID),
            Items:      orderItems,
            Status:     mapOrderStatusToProto(order.Status), // Convert to protobuf enum
            // Include other fields like TotalPrice, etc.
        },
    }

    return response, nil
}

func (s *server) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {
    var order models.Order

    // Start a transaction
    tx := s.db.Begin()

    // Find the order by ID
    if err := tx.First(&order, req.OrderId).Error; err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, status.Errorf(codes.NotFound, "Order with ID '%d' not found", req.OrderId)
        }
        return nil, status.Errorf(codes.Internal, "Error retrieving order: %v", err)
    }

    // Update the order fields
    order.Status = models.OrderStatus(req.Status) // Assuming status is an enum in your protobuf
    // Update other fields as necessary based on req

    // Save the updated order
    if err := tx.Save(&order).Error; err != nil {
        tx.Rollback()
        return nil, status.Errorf(codes.Internal, "Error updating order: %v", err)
    }

    // Commit the transaction
    tx.Commit()

    // Prepare and return the response
    response := &pb.OrderResponse{
        Order: &pb.Order{
            Id:         int64(order.ID),
            CustomerId: int64(order.CustomerID),
            Status:     mapOrderStatusToProto(order.Status),
            // Include other fields and order items if necessary
        },
    }

    return response, nil
}

func (s *server) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
    var orders []models.Order
    query := s.db

    // Implement filtering based on the request, e.g., customer ID
    if req.CustomerId != 0 {
        query = query.Where("customer_id = ?", req.CustomerId)
    }
    // Add more filters as needed based on your application requirements

    // Implement pagination (if required)
    // Example: offset := (req.Page - 1) * req.PageSize
    // query = query.Offset(offset).Limit(req.PageSize)

    // Retrieve the orders from the database
    if err := query.Preload("Items").Find(&orders).Error; err != nil {
        return nil, status.Errorf(codes.Internal, "Error retrieving orders: %v", err)
    }

    // Convert the orders to the protobuf type and create the response
    pbOrders := []*pb.Order{}
    for _, order := range orders {
        pbOrder := &pb.Order{
            Id:         int64(order.ID),
            CustomerId: int64(order.CustomerID),
            Status:     mapOrderStatusToProto(order.Status), // Convert to protobuf enum
            // Map other fields as necessary
        }
        // Map order items if needed
        // ...

        pbOrders = append(pbOrders, pbOrder)
    }

    return &pb.ListOrdersResponse{Orders: pbOrders}, nil
}



func main() {
	db := initDB()
	fmt.Println(db)
    lis, err := net.Listen("tcp", ":50053")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    serv := &server{db: db}
    serv.connectToProductService()
    pb.RegisterOrderServiceServer(s, serv)
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
