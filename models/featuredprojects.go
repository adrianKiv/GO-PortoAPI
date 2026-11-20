package models

import (
	"gorm.io/gorm"
)

type FeaturedProject struct {
	gorm.Model
	Title       string `json:"title" binding:"required"`
	Subtitle    string `json:"subtitle" binding:"required"`
	Detail     string `json:"detail" binding:"required"`
	ImageURL   string `json:"image_url" binding:"required"`
	ImageID     string `json:"image_id"`
	GithubURL   string `json:"github_url" binding:"required"`
	LinkedinURL   string `json:"linkedin_url" binding:"required"`
	YoutubeURL   string `json:"youtube_url" binding:"required"`
}