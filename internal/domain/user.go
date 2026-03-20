package domain

import "gorm.io/gorm"

// User represents a user in the system
type User struct {
	gorm.Model `swaggerignore:"true"`
	Username   string `json:"username" gorm:"unique" extensions:"x-order=0"`
	Password   string `json:"password,omitempty" extensions:"x-order=1"`
}
