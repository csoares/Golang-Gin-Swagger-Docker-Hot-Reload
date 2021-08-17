package model

import "gorm.io/gorm"

// swagger:model
type Evaluation struct {
	gorm.Model `swaggerignore:"true"`
	Rating     int    `json:"Rating" binding:"required,oneof=0 1 2 3 4 5"`
	Note       string `json:"Note"`
}
