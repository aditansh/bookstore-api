package services

import (
	"regexp"
	"time"

	"github.com/aditansh/balkan-task/database"
	"github.com/aditansh/balkan-task/models"
	"github.com/aditansh/balkan-task/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetBookByID(ID uuid.UUID) (models.Book, error) {
	var book models.Book
	result := database.DB.Preload("InCart").Preload("Reviews").Where("id = ?", ID).First(&book)
	if result.Error != nil {
		return models.Book{}, result.Error
	}

	return book, nil
}

func GetBooks() ([]models.Book, error) {
	var books []models.Book
	result := database.DB.Preload("InCart").Preload("Reviews").Find(&books)
	if result.Error != nil {
		return []models.Book{}, result.Error
	}

	return books, nil
}

func GetBooksByVendorID(ID uuid.UUID) ([]models.Book, error) {
	var books []models.Book
	result := database.DB.Preload("InCart").Preload("Reviews").Where("vendor_id = ?", ID).Find(&books)
	if result.Error != nil {
		return []models.Book{}, result.Error
	}

	return books, nil
}

func FindBookByIDAndVendorID(bookID uuid.UUID, vendorID uuid.UUID) (models.Book, error) {
	var book models.Book
	result := database.DB.Preload("InCart").Preload("Reviews").Where("id = ? AND vendor_id = ?", bookID, vendorID).First(&book)
	if result.Error != nil {
		return models.Book{}, result.Error
	}

	return book, nil
}

func CreateBook(payload *schemas.CreateBookSchema, ID uuid.UUID) *fiber.Error {
	books, err := GetBooksByVendorID(ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, book := range books {
		if book.Name == payload.Name {
			return fiber.NewError(fiber.StatusBadRequest, "book with this title already exists")
		}
	}

	book := models.Book{
		Name:        payload.Name,
		Author:      payload.Author,
		Description: payload.Description,
		Categories:  payload.Categories,
		Cost:        payload.Cost,
		Stock:       payload.Stock,
		VendorID:    ID,
	}

	result := database.DB.Create(&book)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func UpdateBook(payload *schemas.UpdateBookSchema, bookID uuid.UUID) *fiber.Error {
	bookToUpdate, err := GetBookByID(bookID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	updates := make(map[string]interface{})
	if payload.Name != "" {
		books, err := GetBooksByVendorID(bookToUpdate.VendorID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		for _, book := range books {
			if book.Name == payload.Name {
				return fiber.NewError(fiber.StatusBadRequest, "Book with this title already exists")
			}
		}
		updates["name"] = payload.Name
	}
	if payload.Author != "" {
		updates["author"] = payload.Author
	}
	if payload.Description != "" {
		updates["description"] = payload.Description
	}
	if payload.Categories != nil {
		updates["categories"] = payload.Categories
	}
	if payload.Cost != 0 {
		updates["cost"] = payload.Cost
	}
	if payload.Stock != 0 {
		updates["stock"] = payload.Stock
	}
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&bookToUpdate).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func DeleteBooks(payload *schemas.DeleteBooksSchema) *fiber.Error {
	var book models.Book
	for _, id := range payload.IDs {
		bookID, err := uuid.Parse(id)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid book ID")
		}

		result := database.DB.Where("id = ?", bookID).Delete(&book)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
		}
	}

	return nil
}

func SearchBooksByAuthor(key string) ([]models.Book, *fiber.Error) {
	allBooks, err := GetBooks()
	if err != nil {
		return []models.Book{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var pattern = `(?i)` + key

	var books []models.Book
	for _, book := range allBooks {
		matched, _ := regexp.Match(pattern, []byte(book.Author))
		if matched {
			books = append(books, book)
		}
	}

	return books, nil
}

func SearchBooksByName(key string) ([]models.Book, *fiber.Error) {
	allBooks, err := GetBooks()
	if err != nil {
		return []models.Book{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var pattern = `(?i)` + key

	var books []models.Book
	for _, book := range allBooks {
		matched, _ := regexp.Match(pattern, []byte(book.Name))
		if matched {
			books = append(books, book)
		}
	}

	return books, nil
}

func FilterBooks(payload *schemas.FilterBooksSchema) ([]models.Book, *fiber.Error) {
	var books []models.Book
	allBooks, err := GetBooks()
	if err != nil {
		return []models.Book{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// if payload.InStock {
	// 	allBooks, err = FilterBooksByStock(allBooks)
	// 	if err != nil {
	// 		return []models.Book{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	// 	}
	// }

	if payload.Categories != nil {
		books, err = FilterBooksByCategory(payload.Categories, allBooks)
		if err != nil {
			return []models.Book{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	if len(books) == 0 {
		books = allBooks
	}

	if payload.PriceRange != nil {
		books, err = FilterBooksByPriceRange(payload.PriceRange[0], payload.PriceRange[1], books)
		if err != nil {
			return []models.Book{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	if len(books) == 0 {
		books = allBooks
	}

	if payload.InStock {
		books, err = FilterBooksByStock(books)
		if err != nil {
			return []models.Book{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return books, nil
}

func FilterBooksByCategory(categories []string, allBooks []models.Book) ([]models.Book, error) {
	var books []models.Book
	for _, book := range allBooks {
		for _, category := range book.Categories {
			for _, c := range categories {
				if category == c {
					books = append(books, book)
				}
			}
		}
	}

	if len(books) == 0 {
		return []models.Book{}, fiber.NewError(fiber.StatusNotFound, "no books found")
	}

	return books, nil
}

func FilterBooksByPriceRange(min float64, max float64, allBooks []models.Book) ([]models.Book, error) {
	var books []models.Book
	for _, book := range allBooks {
		if book.Cost >= min && book.Cost <= max {
			books = append(books, book)
		}
	}

	if len(books) == 0 {
		return []models.Book{}, fiber.NewError(fiber.StatusNotFound, "no books found")
	}

	return books, nil
}

func FilterBooksByStock(allBooks []models.Book) ([]models.Book, error) {
	var books []models.Book
	for _, book := range allBooks {
		if book.Stock > 0 {
			books = append(books, book)
		}
	}

	if len(books) == 0 {
		return []models.Book{}, fiber.NewError(fiber.StatusNotFound, "no books found")
	}

	return books, nil
}
