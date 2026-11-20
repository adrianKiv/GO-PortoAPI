package handler

import (
	"net/http"
	"sync"

	"api/portodata/config"
	"api/portodata/controllers"
	"api/portodata/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	app  *gin.Engine
	once sync.Once
)

func initApp() {
	// Koneksi DB (DATABASE_URI diisi di Vercel Environment Variables)
	config.ConnectDB()

	gin.SetMode(gin.ReleaseMode)
	app = gin.Default()
	app.Use(cors.Default())

	api := app.Group("/api")
	{
		api.POST("/auth/login", controllers.LoginUser)
		api.GET("/ac", controllers.GetACs)
		api.GET("/experiences", controllers.GetExperiences)
		api.GET("/featured-projects", controllers.GetFeaturedProjects)

		protected := api.Group("/")
		protected.Use(middleware.RequireAuth())
		{
			protected.GET("/ac/:id", controllers.GetACByID)
			protected.POST("/ac", controllers.CreateAC)
			protected.PUT("/ac/:id", controllers.UpdateAC)
			protected.DELETE("/ac/:id", controllers.DeleteAC)

			protected.GET("/experiences/:id", controllers.GetExperienceByID)
			protected.POST("/experiences", controllers.CreateExperience)
			protected.PUT("/experiences/:id", controllers.UpdateExperience)
			protected.DELETE("/experiences/:id", controllers.DeleteExperience)

			protected.GET("/featured-projects/:id", controllers.GetFeaturedProjectByID)
			protected.POST("/featured-projects", controllers.CreateFeaturedProject)
			protected.PUT("/featured-projects/:id", controllers.UpdateFeaturedProject)
			protected.DELETE("/featured-projects/:id", controllers.DeleteFeaturedProject)
		}
	}
}

// Export Handler untuk Vercel Serverless Function
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)
	app.ServeHTTP(w, r)
}