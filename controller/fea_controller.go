package controllers

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"os"

	"api/portodata/config"
	"api/portodata/models"
	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

func initImageKit() *imagekit.Client {
	client := imagekit.NewClient(
		option.WithPrivateKey(os.Getenv("IMAGEKIT_PRIVATE_KEY")),
	)
	return &client
}

func CreateFeaturedProject(c *gin.Context) {
	// Ambil file gambar dari form-data dengan key "image"
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gambar wajib diupload",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	// Upload ke ImageKit
	ik := initImageKit()
	uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: header.Filename,
	})
	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to upload image",
			"error":   errUpload.Error(),
		})
		return
	}

	project := models.FeaturedProject{
		Title:       c.PostForm("title"),
		Subtitle:    c.PostForm("subtitle"),
		Detail:     c.PostForm("detail"),
		ImageURL:    uploadRes.URL,
		ImageID:     uploadRes.FileID,
		GithubURL:   c.PostForm("github_url"),
		LinkedinURL: c.PostForm("linkedin_url"),
		YoutubeURL:  c.PostForm("youtube_url"),
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal menyimpan data featured project",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success create featured project",
		"data":    project,
	})
}

func GetFeaturedProjects(c *gin.Context) {
	var projects []models.FeaturedProject

	query := config.DB.Model(&models.FeaturedProject{})

	search := c.Query("search")
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"title ILIKE ? OR subtitle ILIKE ? OR details ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	var totalRows int64
	query.Count(&totalRows)

	pagestr := c.DefaultQuery("page", "1")
	limitstr := c.DefaultQuery("limit", "7")

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

	if err := query.Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal mengambil data featured project",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get all featured projects",
		"data":    projects,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"totalrows":  totalRows,
			"totalpages": totalPage,
		},
	})
}

func GetFeaturedProjectByID(c *gin.Context) {
	var project models.FeaturedProject
	paramID := c.Param("id")

	if err := config.DB.First(&project, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get featured project by id",
		"data":    project,
	})
}

func UpdateFeaturedProject(c *gin.Context) {
	var project models.FeaturedProject
	paramID := c.Param("id")

	if err := config.DB.First(&project, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"error":   err.Error(),
		})
		return
	}

	// Cek apakah user mengirim file gambar baru
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		ik := initImageKit()
		uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
			File:     file,
			FileName: header.Filename,
		})
		if errUpload != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to upload image",
				"error":   errUpload.Error(),
			})
			return
		}

		// Hapus gambar lama di ImageKit jika ada
		if project.ImageID != "" {
			ik.Files.Delete(context.Background(), project.ImageID)
		}

		project.ImageURL = uploadRes.URL
		project.ImageID = uploadRes.FileID
	}

	if title := c.PostForm("title"); title != "" {
		project.Title = title
	}
	if subtitle := c.PostForm("subtitle"); subtitle != "" {
		project.Subtitle = subtitle
	}
	if details := c.PostForm("detail"); details != "" {
		project.Detail += details
	}
	if githubURL := c.PostForm("github_url"); githubURL != "" {
		project.GithubURL = githubURL
	}
	if linkedinURL := c.PostForm("linkedin_url"); linkedinURL != "" {
		project.LinkedinURL = linkedinURL
	}
	if youtubeURL := c.PostForm("youtube_url"); youtubeURL != "" {
		project.YoutubeURL = youtubeURL
	}

	if err := config.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal mengupdate data featured project",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success update featured project",
		"data":    project,
	})
}

func DeleteFeaturedProject(c *gin.Context) {
	var project models.FeaturedProject
	paramID := c.Param("id")

	if err := config.DB.First(&project, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"error":   err.Error(),
		})
		return
	}

	// Hapus file gambar di ImageKit sebelum hapus record database
	if project.ImageID != "" {
		ik := initImageKit()
		ik.Files.Delete(context.Background(), project.ImageID)
	}

	if err := config.DB.Unscoped().Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "gagal menghapus data featured project",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success delete featured project",
	})
}