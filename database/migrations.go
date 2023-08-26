package database

import (
	"github.com/aditansh/balkan-task/models"
	"gorm.io/gorm"
)

func RunMigrations(DB *gorm.DB) {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Vendor{})
	DB.AutoMigrate(&models.Book{})
	DB.AutoMigrate(&models.Review{})
	DB.AutoMigrate(&models.Order{})
	DB.AutoMigrate(&models.OrderItem{})
	DB.AutoMigrate(&models.Cart{})
	DB.AutoMigrate(&models.CartItem{})
}
