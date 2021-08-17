package model

import "gorm.io/gorm"

type Users struct {
	gorm.Model `swaggerignore:"true"`
	Username   string `json:"username" gorm:"unique"`
	Password   string `json:"password,omitempty"`
}
