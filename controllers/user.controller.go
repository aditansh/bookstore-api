package controllers

import (
	"github.com/aditansh/balkan-task/schemas"
	"github.com/aditansh/balkan-task/services"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterUser(c *fiber.Ctx) error {
	var payload schemas.RegisterUserSchema

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

	err := services.RegisterUser(&payload)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  true,
		"message": "User registered successfully",
	})

}

func UpdateUser(c *fiber.Ctx) error {
	var payload schemas.UpdateUserSchema

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

	err := services.UpdateUser(&payload, ID)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	if payload.Email != "" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  true,
			"message": "User updated successfully, please verify your email",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "User updated successfully",
	})
}

func GetUserProfile(c *fiber.Ctx) error {
	ID := c.Locals("ID").(uuid.UUID)

	user, err := services.GetUserByID(ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   user,
	})
}

// admin only routes

func GetAllUsers(c *fiber.Ctx) error {

	users := services.GetAllUsers()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   users,
	})
}

func GetUser(c *fiber.Ctx) error {
	return nil
}

func DeleteUsers(c *fiber.Ctx) error {
	return nil
}

func FlagUser(c *fiber.Ctx) error {
	return nil
}
