package services

import (
	"github.com/aditansh/balkan-task/database"
	"github.com/aditansh/balkan-task/models"
	"github.com/google/uuid"
)

func CreateOrder(ID uuid.UUID, address string) (models.Order, error) {
	newOrder := models.Order{
		UserID:  ID,
		Address: address,
	}

	result := database.DB.Create(&newOrder)
	if result.Error != nil {
		return models.Order{}, result.Error
	}

	return newOrder, nil
}
