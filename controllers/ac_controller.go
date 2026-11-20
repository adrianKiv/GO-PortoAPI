package controllers

import (
	"math"
	"net/http"
	"strconv"

	"api/portodata/config"
	"api/portodata/models"
	"github.com/gin-gonic/gin"
)

func CreateAC(c *gin.Context) {
	var input models.AC

	// Validasi input JSON sesuai tag binding:"required" pada model AC
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "input tidak valid atau ada field yang kosong",
			"error":   err.Error(),
		})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal menyimpan data AC",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success create AC",
		"data":    input,
	})
}

func GetACs(c *gin.Context) {
	var acs []models.AC

	// Inisiasi query
	query := config.DB.Model(&models.AC{})

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

	if err := query.Limit(limit).Offset(offset).Find(&acs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal mengambil data AC",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get all ACs",
		"data":    acs,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalrows":  totalRows,
			"totalpages": totalPage,
		},
	})
}

func GetACByID(c *gin.Context) {
	var ac models.AC
	paramID := c.Param("id")

	if err := config.DB.First(&ac, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get AC by id",
		"data":    ac,
	})
}

func UpdateAC(c *gin.Context) {
	var ac models.AC
	paramID := c.Param("id")

	// Cek apakah data ada di database
	if err := config.DB.First(&ac, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"error":   err.Error(),
		})
		return
	}

	// Struct sementara tanpa binding:"required" agar bisa update sebagian field saja
	var input struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		LinkURL  string `json:"link_url"`
		Date     string `json:"date"`
		Description  string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "format JSON tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if input.Title != "" {
		ac.Title = input.Title
	}
	if input.LinkURL != "" {
		ac.LinkURL = input.LinkURL
	}
	if input.Date != "" {
		ac.Date = input.Date
	}
	if input.Description != "" {
		ac.Description = input.Description
	}

	if err := config.DB.Save(&ac).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal mengupdate data AC",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success update AC",
		"data":    ac,
	})
}

func DeleteAC(c *gin.Context) {
	var ac models.AC
	paramID := c.Param("id")

	if err := config.DB.First(&ac, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"error":   err.Error(),
		})
		return
	}

	// Hard delete mengikuti gaya DeleteEvent (Unscoped)
	if err := config.DB.Unscoped().Delete(&ac).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal menghapus data AC",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success delete AC",
	})
}