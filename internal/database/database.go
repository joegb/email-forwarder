package database

import (
	"fmt"
	"log"
	"time"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	
	"github.com/joegb/email-forwarder/internal/models"
	"github.com/joegb/email-forwarder/internal/config"
)

var DB *gorm.DB

func Connect() {
	cfg := config.GetConfig()
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	
	// 重试连接
	var db *gorm.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Database connection attempt %d failed: %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
	
	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	}
	
	// 自动迁移模型
	if err := db.AutoMigrate(&models.ForwardTarget{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	
	DB = db
	log.Println("Database connected successfully")
}