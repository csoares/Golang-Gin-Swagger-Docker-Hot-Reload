package model

import "github.com/jinzhu/gorm"

type Users struct {
	gorm.Model `swaggerignore:"true"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}
