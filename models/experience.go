package models

import (
	"gorm.io/gorm"
)

type Experience struct {
	gorm.Model
	Title     string `json:"title" binding:"required"`
	Subtitle  string `json:"subtitle" binding:"required"`
	Date      string `json:"date" binding:"required"`
	Detail    string `json:"detail" binding:"required"`
	SortOrder int    `json:"sort_order" gorm:"not null;default:0"`
}
