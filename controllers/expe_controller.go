package controllers

import (
	"math"
	"net/http"
	"strconv"

	"api/portodata/config"
	"api/portodata/models"
	"github.com/gin-gonic/gin"
)

func CreateExperience(c *gin.Context) {
	var input models.Experience

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "input tidak valid atau ada field yang kosong",
			"error":   err.Error(),
		})
		return
	}

	if input.SortOrder < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "sort_order harus lebih besar dari 0",
		})
		return
	}
	if input.SortOrder == 0 {
		sortOrder, err := nextSortOrder(&models.Experience{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "gagal menentukan urutan experience",
				"error":   err.Error(),
			})
			return
		}
		input.SortOrder = sortOrder
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal menyimpan data experience",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success create experience",
		"data":    input,
	})
}

func GetExperiences(c *gin.Context) {
	var experiences []models.Experience

	// Inisiasi query
	query := config.DB.Model(&models.Experience{})

	// Filter pencarian berdasarkan title, subtitle, atau details
	search := c.Query("search")
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"title ILIKE ? OR subtitle ILIKE ? OR details ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Pagination
	var totalRows int64
	query.Count(&totalRows)

	pagestr := c.DefaultQuery("page", "1")
	limitstr := c.DefaultQuery("limit", "6")

	page, errPage := strconv.Atoi(pagestr)
	if errPage != nil || page < 1 {
		page = 1
	}
	limit, errLimit := strconv.Atoi(limitstr)
	if errLimit != nil || limit < 1 {
		limit = 6
	}

	offset := (page - 1) * limit
	totalPage := int(math.Ceil(float64(totalRows) / float64(limit)))

	if err := query.Order("sort_order ASC").Order("id ASC").Limit(limit).Offset(offset).Find(&experiences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal mengambil data experience",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get all experiences",
		"data":    experiences,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalrows":  totalRows,
			"totalpages": totalPage,
		},
	})
}

func GetExperienceByID(c *gin.Context) {
	var experience models.Experience
	paramID := c.Param("id")

	if err := config.DB.First(&experience, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get experience by id",
		"data":    experience,
	})
}

func UpdateExperience(c *gin.Context) {
	var experience models.Experience
	paramID := c.Param("id")

	if err := config.DB.First(&experience, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"error":   err.Error(),
		})
		return
	}

	// Struct sementara tanpa binding:"required" agar bisa update sebagian field saja
	var input struct {
		Title     string `json:"title"`
		Subtitle  string `json:"subtitle"`
		Date      string `json:"date"`
		Detail    string `json:"detail"`
		SortOrder *int   `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "format JSON tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if input.Title != "" {
		experience.Title = input.Title
	}
	if input.Subtitle != "" {
		experience.Subtitle = input.Subtitle
	}
	if input.Date != "" {
		experience.Date = input.Date
	}
	if input.Detail != "" {
		experience.Detail = input.Detail
	}
	if input.SortOrder != nil {
		if *input.SortOrder < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "sort_order harus lebih besar dari 0",
			})
			return
		}
		experience.SortOrder = *input.SortOrder
	}

	if err := config.DB.Save(&experience).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal mengupdate data experience",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success update experience",
		"data":    experience,
	})
}

func DeleteExperience(c *gin.Context) {
	var experience models.Experience
	paramID := c.Param("id")

	if err := config.DB.First(&experience, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"error":   err.Error(),
		})
		return
	}

	if err := config.DB.Unscoped().Delete(&experience).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal menghapus data experience",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success delete experience",
	})
}
