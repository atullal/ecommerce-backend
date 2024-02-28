package main

import (
	"os"
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "github.com/atullal/ecommerce-backend-protobuf/user"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "user-service/models"
    "fmt"
    "golang.org/x/crypto/bcrypt"
    "time"
    "math"
    "utils/jwt/jwt.go"
)

type server struct {
    pb.UserServiceServer
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
    if err := db.AutoMigrate(&models.User{}); err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }
    fmt.Println("Database connection successful")
    return db
}

func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.UserResponse, error) {
    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    user := models.User{Username: in.Username, Email: in.Email, Password: string(hashedPassword)}
    result := s.db.Create(&user)
    if result.Error != nil {
        fmt.Println("Error: ", result.Error)
        return nil, result.Error
    }

    // Generate JWT token
    token, err := GenerateToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &pb.UserResponse{Id: fmt.Sprintf("%d", user.ID), Username: user.Username, Email: user.Email, Token: token}, nil
}


func (s *server) AuthenticateUser(ctx context.Context, in *pb.AuthenticateUserRequest) (*pb.UserResponse, error) {
    var user models.User
    if err := s.db.Where("email = ?", in.Email).First(&user).Error; err != nil {
        return nil, err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)); err != nil {
        return nil, err
    }

    // Generate JWT token
    token, err := GenerateToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &pb.UserResponse{Id: fmt.Sprintf("%d", user.ID), Username: user.Username, Email: user.Email, Token: token}, nil
}


func main() {
	db := initDB()
	fmt.Println(db)
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    serv := &server{db: db}
    pb.RegisterUserServiceServer(s, serv)
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
