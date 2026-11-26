package config

import (
	"log"
	"os"

	"api/portodata/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URI")
	if dsn == "" {
		log.Fatal("DATABASE_URI is not set")
	}

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt: false,
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = database.AutoMigrate(&models.User{}, &models.Experience{}, &models.AC{}, &models.FeaturedProject{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	for _, model := range []any{&models.Experience{}, &models.AC{}, &models.FeaturedProject{}} {
		if err := database.Model(model).Where("sort_order = 0").Update("sort_order", gorm.Expr("id")).Error; err != nil {
			log.Fatal("Failed to initialize data ordering:", err)
		}
	}

	DB = database
	log.Println("Database connection established")
}
