package models

import (
    "gorm.io/gorm"
)

// Order represents the main structure for an order
type Order struct {
    gorm.Model
    CustomerID  uint        `gorm:"index"`  // Assuming you have customer IDs as uint
    Items       []OrderItem // Association with OrderItem
    Status      OrderStatus // Custom type defined below
    TotalPrice  float64     // Total price of the order
    // Add other fields like shipping address, payment details, etc.
}

// OrderItem represents an item in an order
type OrderItem struct {
    gorm.Model
    OrderID   uint    // Foreign key for the Order
    ProductID uint    // Assuming product IDs as uint
    Quantity  int     // Quantity of the product
    Price     float64 // Price of an individual item
    Version   int     // Optimistic locking version
    // You can add more fields if necessary
}

// OrderStatus represents the status of an order
type OrderStatus string

// Enum values for OrderStatus
const (
    StatusPending    OrderStatus = "PENDING"
    StatusConfirmed  OrderStatus = "CONFIRMED"
    StatusShipped    OrderStatus = "SHIPPED"
    StatusDelivered  OrderStatus = "DELIVERED"
    StatusCancelled  OrderStatus = "CANCELLED"
)
