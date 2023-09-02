package controllers

import (
	"github.com/aditansh/balkan-task/schemas"
	"github.com/aditansh/balkan-task/services"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetCart(c *fiber.Ctx) error {
	ID := c.Locals("ID").(uuid.UUID)

	cart, err := services.GetCart(ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   cart,
	})
}

func GetCartValue(c *fiber.Ctx) error {
	ID := c.Locals("ID").(uuid.UUID)

	value, err := services.GetCartValue(ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   value,
	})
}

func AddToCart(c *fiber.Ctx) error {
	var payload schemas.ModifyCartSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid payload",
		})
	}

	errors := utils.ValidateStruct(&payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": errors,
		})
	}

	ID := c.Locals("ID").(uuid.UUID)
	err := services.AddToCart(ID, &payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "added to cart",
	})
}

func UpdateCart(c *fiber.Ctx) error {
	var payload schemas.ModifyCartSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid payload",
		})
	}

	errors := utils.ValidateStruct(&payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": errors,
		})
	}

	ID := c.Locals("ID").(uuid.UUID)
	err := services.UpdateCart(ID, &payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "updated cart",
	})
}

func ClearCart(c *fiber.Ctx) error {
	ID := c.Locals("ID").(uuid.UUID)
	err := services.ClearCart(ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "cleared cart",
	})
}

func RemoveFromCart(c *fiber.Ctx) error {
	var payload schemas.RemoveFromCart
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid payload",
		})
	}

	errors := utils.ValidateStruct(&payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": errors,
		})
	}

	ID := c.Locals("ID").(uuid.UUID)
	err := services.RemoveFromCart(ID, &payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "removed from cart",
	})
}

func Checkout(c *fiber.Ctx) error {
	var payload schemas.CheckoutSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid payload",
		})
	}

	errors := utils.ValidateStruct(&payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": errors,
		})
	}

	ID := c.Locals("ID").(uuid.UUID)
	orderID, err := services.Checkout(ID, &payload)
	if err != nil {
		result := services.ClearOrderItems(orderID)
		if result != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "checked out",
	})
}

// admin only routes

func GetAllCarts(c *fiber.Ctx) error {
	carts, err := services.GetAllCarts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   carts,
	})

}

func GetUserCart(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid user id",
		})
	}

	cart, err := services.GetCart(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   cart,
	})
}

func ModifyUserCart(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid user id",
		})
	}

	var payload schemas.ModifyCartSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid payload",
		})
	}

	errors := utils.ValidateStruct(&payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": errors,
		})
	}

	errr := services.UpdateCart(userID, &payload)
	if errr != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "updated cart",
	})
}

func ClearUserCart(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid user id",
		})
	}

	errr := services.ClearCart(userID)
	if errr != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "cleared cart",
	})
}
