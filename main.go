package main

import (
	"api/portodata/controllers"
	"api/portodata/config"
	"api/portodata/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)

func main() {

	_ = godotenv.Load()

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }	

	config.ConnectDB()

	router := gin.Default()
	router.Use(cors.Default())

	api := router.Group("/api")
	{

		api.GET("/ac", controllers.GetACs)
		api.GET("/experiences", controllers.GetExperiences)
		api.GET("/featured-projects", controllers.GetFeaturedProjects)
		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)
		
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


	router.Run()
}
