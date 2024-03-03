// models/user.go
package models

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username  string
    Email     string `gorm:"unique"`
    Password  string
    Role	  int
}
