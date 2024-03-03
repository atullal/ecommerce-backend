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
    pb "github.com/atullal/ecommerce-backend-protobuf/product"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "product-service/models"
    "fmt"
    "time"
    "math"
)

type server struct {
    pb.ProductServiceServer
    db *gorm.DB
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
    if err := db.AutoMigrate(&models.Product{}); err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }
    fmt.Println("Database connection successful")
    return db
}

// AddProduct handles the creation of a new product
func (s *server) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.ProductResponse, error) {
    // Create a new Product models instance from the request
    newProduct := models.Product{
        Name:        req.Name,
        Description: req.Description,
        Price:       float64(req.Price),      // Convert to float64
        Quantity:    int(req.Quantity),       // Convert to int
    }

    // Save the new product to the database
    if err := s.db.Create(&newProduct).Error; err != nil {
        return nil, err // Handle and return the error appropriately
    }

    // Prepare and return the response
    response := &pb.ProductResponse{
	    Product: &pb.Product{
	        Id:          int64(newProduct.ID),  // Convert to int64
	        Name:        newProduct.Name,
	        Description: newProduct.Description,
	        Price:       float32(newProduct.Price), // Convert to float32
	        Quantity:    int32(newProduct.Quantity),
			Version:    int64(newProduct.Version),
	    },
    }

    return response, nil
}

// GetProduct handles fetching a product by ID
func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
    var product models.Product

    // Retrieve the product by ID from the database
    result := s.db.First(&product, req.Id)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, status.Errorf(codes.NotFound, "Product with ID '%d' not found", req.Id)
        }
        return nil, status.Errorf(codes.Internal, "Error retrieving product: %v", result.Error)
    }

    // Prepare and return the response
    response := &pb.ProductResponse{
        Product: &pb.Product{
            Id:          int64(product.ID),
            Name:        product.Name,
            Description: product.Description,
            Price:       float32(product.Price),
            Quantity:    int32(product.Quantity),
            Version:	 int64(product.Version),
        },
    }

    return response, nil
}

func (s *server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
    var product models.Product

    // Start a transaction
    tx := s.db.Begin()

    // Find the product by ID
    if err := tx.First(&product, req.Id).Error; err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, status.Errorf(codes.NotFound, "Product with ID '%d' not found", req.Id)
        }
        return nil, status.Errorf(codes.Internal, "Error retrieving product: %v", err)
    }

    // Check if the version matches
    if int64(req.Version) != int64(product.Version) {
        tx.Rollback()
        return nil, status.Errorf(codes.Aborted, "Update aborted due to version mismatch")
    }

    // Update the product fields
    product.Name = req.Name
    product.Description = req.Description
    product.Price = float64(req.Price)
    product.Version++ // Increment the version

    // Save the updated product
    if err := tx.Save(&product).Error; err != nil {
        tx.Rollback()
        return nil, status.Errorf(codes.Internal, "Error updating product: %v", err)
    }

    // Commit the transaction
    tx.Commit()

    // Prepare and return the response
    return &pb.ProductResponse{
        Product: &pb.Product{
            Id:          int64(product.ID),
            Name:        product.Name,
            Description: product.Description,
            Price:       float32(product.Price),
            Quantity:    int32(product.Quantity),
            Version:     int64(product.Version),
        },
    }, nil
}

func (s *server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
    // Find the product by ID
    var product models.Product
    result := s.db.First(&product, req.Id)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, status.Errorf(codes.NotFound, "Product with ID '%d' not found", req.Id)
        }
        return nil, status.Errorf(codes.Internal, "Error retrieving product: %v", result.Error)
    }

    // Delete the product
    if err := s.db.Delete(&product).Error; err != nil {
        return nil, status.Errorf(codes.Internal, "Error deleting product: %v", err)
    }

    // Return a success response
    return &pb.DeleteProductResponse{Success: true}, nil
}

func (s *server) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
    var products []models.Product
    query := s.db

    // Implement filtering based on the request, e.g., search keyword, categories
    if req.SearchKeyword != "" {
        query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+req.SearchKeyword+"%", "%"+req.SearchKeyword+"%")
    }
    if len(req.Categories) > 0 {
        query = query.Where("category IN ?", req.Categories)
    }
    // Add more filters as needed

    // Implement pagination
    offset := (req.Page - 1) * req.PageSize
    query = query.Offset(int(offset)).Limit(int(req.PageSize))

    // Retrieve the products from the database
    if err := query.Find(&products).Error; err != nil {
        return nil, status.Errorf(codes.Internal, "Error retrieving products: %v", err)
    }

    // Convert the products to the protobuf type and create the response
    pbProducts := []*pb.Product{}
    for _, product := range products {
        pbProducts = append(pbProducts, &pb.Product{
            Id:          int64(product.ID),
            Name:        product.Name,
            Description: product.Description,
            Price:       float32(product.Price),
            Quantity:    int32(product.Quantity),
            Version:    int64(product.Version),
            // Map other fields as necessary
        })
    }

    return &pb.ListProductsResponse{Products: pbProducts}, nil
}

func (s *server) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryRequest) (*pb.InventoryResponse, error) {
    var product models.Product

    // Start a transaction
    tx := s.db.Begin()

    // Retrieve the product with a SELECT FOR UPDATE lock
    if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, req.ProductId).Error; err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, status.Errorf(codes.NotFound, "Product with ID '%d' not found", req.ProductId)
        }
        return nil, status.Errorf(codes.Internal, "Error retrieving product: %v", err)
    }

    // Check the version for optimistic locking
    if int64(req.Version) != int64(product.Version) {
        tx.Rollback()
        return nil, status.Errorf(codes.Aborted, "Inventory update aborted due to version mismatch")
    }

    // Update the inventory
    product.Quantity += int(req.QuantityChange) // Adjust this line based on your business logic
    product.Version++

    if product.Quantity < 0 {
    	tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "Insufficient inventory")
    }

    // Save the changes
    if err := tx.Save(&product).Error; err != nil {
        tx.Rollback()
        return nil, status.Errorf(codes.Internal, "Error updating inventory: %v", err)
    }

    // Commit the transaction
    tx.Commit()

    // Prepare and return the response
    return &pb.InventoryResponse{
        ProductId: int64(product.ID),
        Quantity:  int32(product.Quantity),
        Version:   int64(product.Version),
    }, nil
}

func (s *server) UpdateMultipleInventories(ctx context.Context, req *pb.UpdateMultipleInventoriesRequest) (*pb.InventoriesResponse, error) {
    fmt.Println("UpdateMultipleInventories: ", req)
	res := &pb.InventoriesResponse{}
    res.Inventories = make([]*pb.InventoryResponse, 0)
	// Start a transaction
    tx := s.db.Begin()

	for _, inventory := range req.InventoryUpdates {
		if inventory.ProductId == 0 {
			continue
		}
		var product models.Product

	    // Retrieve the product with a SELECT FOR UPDATE lock
	    if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, inventory.ProductId).Error; err != nil {
	        tx.Rollback()
	        if errors.Is(err, gorm.ErrRecordNotFound) {
	            return nil, status.Errorf(codes.NotFound, "Product with ID '%d' not found", inventory.ProductId)
	        }
	        return nil, status.Errorf(codes.Internal, "Error retrieving product: %v", err)
	    }
		fmt.Println("product: ", product)
		fmt.Println("inventory: ", inventory)
	    // Check the version for optimistic locking
	    if int64(inventory.Version) != int64(product.Version) {
	        tx.Rollback()
	        return nil, status.Errorf(codes.Aborted, "Inventory update aborted due to version mismatch")
	    }

	    // Update the inventory
	    product.Quantity += int(inventory.QuantityChange) // Adjust this line based on your business logic
	    product.Version++

	    if product.Quantity < 0 {
	    	tx.Rollback()
				return nil, status.Errorf(codes.InvalidArgument, "Insufficient inventory")
	    }

	    // Save the changes
	    if err := tx.Save(&product).Error; err != nil {
	        tx.Rollback()
	        return nil, status.Errorf(codes.Internal, "Error updating inventory: %v", err)
	    }
		res.Inventories = append(res.Inventories, &pb.InventoryResponse{
			ProductId: int64(product.ID),
	        Quantity:  int32(product.Quantity),
	        Version:   int64(product.Version),
		})
    }

    // Commit the transaction
    tx.Commit()

    // Prepare and return the response
    return res, nil
}

func (s *server) GetInventory(ctx context.Context, req *pb.GetInventoryRequest) (*pb.InventoryResponse, error) {
    var product models.Product

    // Retrieve the product by ID from the database
    result := s.db.First(&product, req.ProductId)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, status.Errorf(codes.NotFound, "Product with ID '%d' not found", req.ProductId)
        }
        return nil, status.Errorf(codes.Internal, "Error retrieving product: %v", result.Error)
    }

    // Prepare and return the response with the inventory details
    response := &pb.InventoryResponse{
        ProductId: int64(product.ID),
        Quantity:  int32(product.Quantity),
        Version:   int64(product.Version),
    }

    return response, nil
}


func main() {
	db := initDB()
	fmt.Println(db)
    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    serv := &server{db: db}
    pb.RegisterProductServiceServer(s, serv)
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
