package main

import (
    "log"
    "github.com/gin-gonic/gin"
    pb "github.com/atullal/ecommerce-backend-protobuf/user" // Replace with the correct import path
    "rest-service/handlers"
    "rest-service/services"
    "google.golang.org/grpc"
    "fmt"
)

type server struct {
	RestServer *gin.Engine
	UserServiceConnection *grpc.ClientConn
}

func (s *server) AddUserRoutes(userHandler handlers.UserHandler) {
	// Set up routes
    s.RestServer.POST("/user", userHandler.CreateUser)
    s.RestServer.POST("/user/authenticate", userHandler.AuthenticateUser)
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
    userService := services.NewUserService(pb.NewUserServiceClient(userServiceConnection))
    userHandler := handlers.UserHandler{UserService: userService}

    s.AddUserRoutes(userHandler)
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


    // Start the server
    if err := s.RestServer.Run(); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
