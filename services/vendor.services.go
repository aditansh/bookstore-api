package services

import (
	"github.com/aditansh/balkan-task/database"
	"github.com/aditansh/balkan-task/models"
	"github.com/google/uuid"
)

func GetVendorByID(id uuid.UUID) (models.Vendor, error) {
	var vendor models.Vendor
	result := database.DB.Where("id = ?", id).First(&vendor)
	if result.Error != nil {
		return models.Vendor{}, result.Error
	}

	return vendor, nil
}
