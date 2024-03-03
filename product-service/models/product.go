package models

import (
    "gorm.io/gorm"
)

type Product struct {
    gorm.Model
    Name        string
    Description string
    Price       float64
    Quantity    int
    Version     int // Optimistic locking version
}
