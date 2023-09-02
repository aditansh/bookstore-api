package services

import (
	"github.com/aditansh/balkan-task/database"
	"github.com/aditansh/balkan-task/models"
	"github.com/google/uuid"
)

func GetAllOrders() ([]models.Order, error) {
	var orders []models.Order

	result := database.DB.Preload("Items").Find(&orders)
	if result.Error != nil {
		return []models.Order{}, result.Error
	}

	return orders, nil
}

func GetOrders(ID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order

	result := database.DB.Preload("Items").Where("user_id = ?", ID).Find(&orders)
	if result.Error != nil {
		return []models.Order{}, result.Error
	}

	return orders, nil
}

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

func GetOrder(ID uuid.UUID) (models.Order, error) {
	var order models.Order

	result := database.DB.Preload("Items").Where("id = ?", ID).First(&order)
	if result.Error != nil {
		return models.Order{}, result.Error
	}

	return order, nil
}

func GetOrderByUserIDAndOrderID(userID uuid.UUID, orderID uuid.UUID) (models.Order, error) {
	var order models.Order

	result := database.DB.Preload("Items").Where("id = ? AND user_id = ?", orderID, userID).First(&order)
	if result.Error != nil {
		return models.Order{}, result.Error
	}

	return order, nil
}

func ClearOrderItems(ID string) error {
	orderID, err := uuid.Parse(ID)
	if err != nil {
		return err
	}

	order, err := GetOrder(orderID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		book, err := GetBookByID(item.BookID)
		if err != nil {
			return err
		}

		updates := make(map[string]interface{})
		updates["stock"] = book.Stock + item.Quantity

		result := database.DB.Model(&book).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
	}

	result := database.DB.Where("id = ?", orderID).Delete(&order)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func SearchOrdersByAuthor(author string, userID uuid.UUID) ([]models.Order, error) {
	orders, err := GetOrders(userID)
	if err != nil {
		return []models.Order{}, err
	}

	var result []models.Order

	for _, order := range orders {
		for _, item := range order.Items {
			if item.Author == author {
				result = append(result, order)
				break
			}
		}
	}

	return result, nil
}

func SearchOrdersByBookName(bookName string, userID uuid.UUID) ([]models.Order, error) {
	orders, err := GetOrders(userID)
	if err != nil {
		return []models.Order{}, err
	}

	var result []models.Order

	for _, order := range orders {
		for _, item := range order.Items {
			if item.BookName == bookName {
				result = append(result, order)
				break
			}
		}
	}

	return result, nil
}

func SearchOrdersByVendorName(vendorName string, userID uuid.UUID) ([]models.Order, error) {
	orders, err := GetOrders(userID)
	if err != nil {
		return []models.Order{}, err
	}

	var result []models.Order

	for _, order := range orders {
		for _, item := range order.Items {
			if item.VendorName == vendorName {
				result = append(result, order)
				break
			}
		}
	}

	return result, nil
}

func SearchAllOrdersByAuthor(author string) ([]models.Order, error) {
	orders, err := GetAllOrders()
	if err != nil {
		return []models.Order{}, err
	}

	var result []models.Order

	for _, order := range orders {
		for _, item := range order.Items {
			if item.Author == author {
				result = append(result, order)
				break
			}
		}
	}

	return result, nil
}

func SearchAllOrdersByBookName(bookName string) ([]models.Order, error) {
	orders, err := GetAllOrders()
	if err != nil {
		return []models.Order{}, err
	}

	var result []models.Order

	for _, order := range orders {
		for _, item := range order.Items {
			if item.BookName == bookName {
				result = append(result, order)
				break
			}
		}
	}

	return result, nil
}

func SearchAllOrdersByVendorName(vendorName string) ([]models.Order, error) {
	orders, err := GetAllOrders()
	if err != nil {
		return []models.Order{}, err
	}

	var result []models.Order

	for _, order := range orders {
		for _, item := range order.Items {
			if item.VendorName == vendorName {
				result = append(result, order)
				break
			}
		}
	}

	return result, nil
}
