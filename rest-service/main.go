package main

import (
    "log"
    "github.com/gin-gonic/gin"
    pb "github.com/atullal/ecommerce-backend/user-service/gen" // Replace with the correct import path
    "rest-service/handlers"
    "rest-service/services"
    "google.golang.org/grpc"
)

func main() {
    // Initialize the REST server
    r := gin.Default()

    // Initialize user-related components
    initializeUserComponents(r)

    // Start the server
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}

// initializeUserComponents sets up everything related to user handling
func initializeUserComponents(r *gin.Engine) {
    // Set up a connection to the gRPC server.
    conn, err := grpc.Dial("user-service:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }
    defer conn.Close()

    // Initialize the UserService and UserHandler
    userService := services.NewUserService(pb.NewUserServiceClient(conn))
    userHandler := handlers.UserHandler{UserService: userService}

    // Set up routes
    r.POST("/user", userHandler.CreateUser)
    r.POST("/user/authenticate", userHandler.AuthenticateUser)

    // Additional user-related routes can be added here
}
