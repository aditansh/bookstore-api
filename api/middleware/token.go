package middleware

import (
	"github.com/aditansh/balkan-task/services"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
)

// VerifyAdminToken verifies the admin token
func VerifyAdminToken(c *fiber.Ctx) error {
	email, err := utils.VerifyToken(c)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	user, errr := services.GetUserByEmail(email)
	if errr != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": errr.Error(),
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Not an admin",
		})
	}

	c.Locals("ID", user.ID)

	return c.Next()
}

// VerifyUserToken verifies the user token
func VerifyUserToken(c *fiber.Ctx) error {
	email, err := utils.VerifyToken(c)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	user, errr := services.GetUserByEmail(email)
	if errr != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": errr.Error(),
		})
	}

	if user.Role != "user" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Not an user",
		})
	}

	c.Locals("ID", user.ID)

	return c.Next()
}

// VerifyToken verifies the token for all
func VerifyToken(c *fiber.Ctx) error {

	email, err := utils.VerifyToken(c)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	c.Locals("Email", email)

	return c.Next()
}

// VerifyVendorToken verifies the vendor token
func VerifyVendorToken(c *fiber.Ctx) error {
	email, err := utils.VerifyToken(c)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	vendor, errr := services.GetVendorByEmail(email)
	if errr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	c.Locals("ID", vendor.ID)

	return c.Next()
}
