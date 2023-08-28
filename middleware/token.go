package middleware

import (
	"github.com/aditansh/balkan-task/services"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// VerifyAdminToken verifies the admin token
func VerifyAdminToken(c *fiber.Ctx) error {
	user, err := utils.VerifyToken(c, services.GetUserByEmail)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
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

// VerifyVendorToken verifies the vendor token
func VerifyVendorToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "No authorization header found",
		})
	}
	token := utils.GetToken(authHeader)

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Token not found",
		})
	}

	res, err := utils.ValidateToken(token, viper.GetString("ACCESS_TOKEN_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	_, err = services.GetVendorByID(res.ID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	c.Locals("ID", res.ID)

	return c.Next()
}

// VerifyUserToken verifies the user token
func VerifyUserToken(c *fiber.Ctx) error {
	user, err := utils.VerifyToken(c, services.GetUserByEmail)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
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

	user, err := utils.VerifyToken(c, services.GetUserByEmail)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	c.Locals("Email", user.Email)

	return c.Next()
}
