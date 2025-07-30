package db

import (
	"fmt"
	"log"

	"trading_api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=db user=tradingbotuser password=securepassword dbname=tradingbotdb port=5432 sslmode=disable"

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to DB:", err)
	}

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("❌ Failed to auto-migrate:", err)
	}

	fmt.Println("✅ PostgreSQL connected and schema migrated")
}
