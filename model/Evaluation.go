package model

import "github.com/jinzhu/gorm"

// swagger:model
type Evaluation struct {
	gorm.Model `swaggerignore:"true"`
	Rating     int    `json:"Rating"`
	Note       string `json:"Note"`
}
