package models

import (
	"gorm.io/gorm"
)

type AC struct {
	gorm.Model
	Type        string `json:"type" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Date        string `json:"date" binding:"required"`
	Description string `json:"description" binding:"required"`
	LinkURL     string `json:"link_url"`
	SortOrder   int    `json:"sort_order" gorm:"not null;default:0"`
}
