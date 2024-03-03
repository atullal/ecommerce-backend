package main

import (
    "log"
    "github.com/gin-gonic/gin"
    userpb "github.com/atullal/ecommerce-backend-protobuf/user"
    productpb "github.com/atullal/ecommerce-backend-protobuf/product"
    orderpb "github.com/atullal/ecommerce-backend-protobuf/order"
    "rest-service/handlers"
    "rest-service/services"
    "rest-service/middlewares"
    "google.golang.org/grpc"
    "fmt"
)

type server struct {
	RestServer *gin.Engine
	UserServiceConnection *grpc.ClientConn
	ProductServiceConnection *grpc.ClientConn
	OrderServiceConnection *grpc.ClientConn
}

func (s *server) AddUserRoutes(userHandler handlers.UserHandler) {
	// Set up routes
    s.RestServer.POST("/user", userHandler.CreateUser)
    s.RestServer.POST("/user/authenticate", userHandler.AuthenticateUser)
}

func (s *server) AddProductRoutes(productHandler handlers.ProductHandler) {
    // Set up product-related routes
    s.RestServer.GET("/products", productHandler.ListProducts)
    s.RestServer.GET("/product/:id", productHandler.GetProduct)
    // Add any other product routes here
    admin := s.RestServer.Group("/")
    admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
    {
        admin.DELETE("/product/:id", productHandler.DeleteProduct)
        admin.POST("/product", productHandler.AddProduct)
        admin.PUT("/product/:id", productHandler.UpdateProduct)
        admin.PUT("/product/inventory/:id", productHandler.UpdateInventory)
        admin.GET("/product/inventory/:id", productHandler.GetInventory)
    }
}

func (s *server) AddOrderRoutes(orderHandler handlers.OrderHandler) {
    // Apply middleware
    authenticated := s.RestServer.Group("/")
    authenticated.Use(middleware.AuthMiddleware())
    {
        authenticated.POST("/order", orderHandler.CreateOrder)
        authenticated.GET("/order/:id", orderHandler.GetOrder)
        authenticated.GET("/orders", orderHandler.ListOrders)
    }

    // Add any other product routes here
    admin := s.RestServer.Group("/")
    admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
    {
        admin.PUT("/order/:id", orderHandler.UpdateOrder)
    }
}

// initializeUserComponents sets up everything related to user handling
func (s *server) InitializeUserComponents() {
    // Set up a connection to the gRPC server.
    userServiceConnection, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }
    fmt.Println("Connected to gRPC server")

    // Initialize the UserService and UserHandler
    userService := services.NewUserService(userpb.NewUserServiceClient(userServiceConnection))
    userHandler := handlers.UserHandler{UserService: userService}

    s.AddUserRoutes(userHandler)
}

// initializeUserComponents sets up everything related to user handling
func (s *server) InitializeProductComponents() {
    // Set up a connection to the gRPC server.
    productServiceConnection, err := grpc.Dial("0.0.0.0:50052", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }
    fmt.Println("Connected to gRPC server")

    // Initialize the ProductService and ProductHandler
    productService := services.NewProductService(productpb.NewProductServiceClient(productServiceConnection))
    productHandler := handlers.ProductHandler{ProductService: productService}

    s.AddProductRoutes(productHandler)
}

// initializeUserComponents sets up everything related to order handling
func (s *server) InitializeOrderComponents() {
    // Set up a connection to the gRPC server.
    orderServiceConnection, err := grpc.Dial("0.0.0.0:50053", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }
    fmt.Println("Connected to gRPC server")

    orderService := services.NewOrderService(orderpb.NewOrderServiceClient(orderServiceConnection))
    orderHandler := handlers.OrderHandler{OrderService: orderService}

    s.AddOrderRoutes(orderHandler)
}

func (s *server) InitializeRestService() {
	// Set up routes
	s.RestServer.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})
}

func (s *server) Close() {
	s.UserServiceConnection.Close()
	s.ProductServiceConnection.Close()
	s.RestServer = nil
}

func main() {
	s := server{}
	defer s.Close()
    // Initialize the REST server
    s.RestServer = gin.Default()
    // Initialize rest components
    s.InitializeRestService()

    // Initialize user-related components
    s.InitializeUserComponents()

    // Initialize product-related components
    s.InitializeProductComponents()

    // Initialize order-related components
    s.InitializeOrderComponents()

    // Start the server
    if err := s.RestServer.Run(); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
