package controllers

import (
	"github.com/aditansh/balkan-task/schemas"
	"github.com/aditansh/balkan-task/services"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetBooks(c *fiber.Ctx) error {
	books, err := services.GetBooks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   books,
	})
}

func GetBook(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid book ID",
		})
	}

	book, err := services.GetBookByID(ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   book,
	})
}

func SearchBooks(c *fiber.Ctx) error {
	var payload schemas.SearchBooksSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"errors": errors,
		})
	}

	if payload.SearchBy == "name" {
		books, err := services.SearchBooksByName(payload.Search)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"data":   books,
		})
	} else if payload.SearchBy == "author" {
		books, err := services.SearchBooksByAuthor(payload.Search)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"data":   books,
		})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid search parameter",
		})
	}
}

func FilterBooks(c *fiber.Ctx) error {
	var payload schemas.FilterBooksSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"errors": errors,
		})
	}

	if payload.PriceRange != nil {
		err := utils.ValidatePriceRange(payload.PriceRange)
		if !err {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Invalid Price Range",
			})
		}
	}

	books, errr := services.FilterBooks(&payload)
	if errr != nil {
		return c.Status(errr.Code).JSON(fiber.Map{
			"status":  false,
			"message": errr.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   books,
	})
}

// vendor routes

func CreateBook(c *fiber.Ctx) error {
	var payload schemas.CreateBookSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"errors": errors,
		})
	}

	ID := c.Locals("ID").(uuid.UUID)

	err := services.CreateBook(&payload, ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  true,
		"message": "Book created successfully",
	})
}

func UpdateBook(c *fiber.Ctx) error {
	var payload schemas.UpdateBookSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"errors": errors,
		})
	}

	bookIdStr := c.Params("id")
	bookID, err := uuid.Parse(bookIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid book ID",
		})
	}
	errr := services.UpdateBook(&payload, bookID)
	if errr != nil {
		return c.Status(errr.Code).JSON(fiber.Map{
			"status":  false,
			"message": errr.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Book updated successfully",
	})
}

func DeleteBooks(c *fiber.Ctx) error {
	var payload schemas.DeleteBooksSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"errors": errors,
		})
	}

	err := services.DeleteBooks(&payload)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Books deleted successfully",
	})
}
