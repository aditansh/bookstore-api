package controllers

import (
	"github.com/aditansh/balkan-task/schemas"
	"github.com/aditansh/balkan-task/services"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetOrders(c *fiber.Ctx) error {
	ID := c.Locals("ID").(uuid.UUID)

	orders, err := services.GetOrders(ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   orders,
	})
}

func GetOrder(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	userID := c.Locals("ID").(uuid.UUID)
	order, err := services.GetOrderByUserIDAndOrderID(userID, orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   order,
	})
}

func SearchUserOrders(c *fiber.Ctx) error {
	var payload schemas.SearchOrdersSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid payload",
		})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": errors,
		})
	}

	var userID = c.Locals("ID").(uuid.UUID)

	if payload.SearchBy == "bookName" {
		orders, err := services.SearchOrdersByBookName(payload.Query, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"data":   orders,
		})
	} else if payload.SearchBy == "author" {
		orders, err := services.SearchOrdersByAuthor(payload.Query, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"data":   orders,
		})
	} else if payload.SearchBy == "vendorName" {
		orders, err := services.SearchOrdersByVendorName(payload.Query, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"data":   orders,
		})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid searchBy field",
		})
	}
}

// func FilterUserOrders(c *fiber.Ctx) error {
// 	return nil
// }

// admmin only routes

func GetAllOrders(c *fiber.Ctx) error {
	orders, err := services.GetAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   orders,
	})
}

func GetUserOrders(c *fiber.Ctx) error {
	username := (c.Params("username"))

	user, err := services.GetUserByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	orders, err := services.GetOrders(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   orders,
	})
}

func GetUserOrder(c *fiber.Ctx) error {
	username := (c.Params("username"))

	user, err := services.GetUserByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	orderID, err := uuid.Parse(c.Params("orderid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	order, err := services.GetOrderByUserIDAndOrderID(user.ID, orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   order,
	})
}

func SearchOrders(c *fiber.Ctx) error {
	var payload schemas.SearchOrdersSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid payload",
		})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": errors,
		})
	}

	username := c.Params("username")

	if username == ":username" {
		if payload.SearchBy == "bookName" {
			orders, err := services.SearchAllOrdersByBookName(payload.Query)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"data":   orders,
			})
		} else if payload.SearchBy == "author" {
			orders, err := services.SearchAllOrdersByAuthor(payload.Query)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"data":   orders,
			})
		} else if payload.SearchBy == "vendorName" {
			orders, err := services.SearchAllOrdersByVendorName(payload.Query)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"data":   orders,
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "invalid searchBy field",
			})
		}
	} else {
		user, err := services.GetUserByUsername(username)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		if payload.SearchBy == "bookName" {
			orders, err := services.SearchOrdersByBookName(payload.Query, user.ID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"data":   orders,
			})
		} else if payload.SearchBy == "author" {
			orders, err := services.SearchOrdersByAuthor(payload.Query, user.ID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"data":   orders,
			})
		} else if payload.SearchBy == "vendorName" {
			orders, err := services.SearchOrdersByVendorName(payload.Query, user.ID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"data":   orders,
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "invalid searchBy field",
			})
		}
	}
}

// func FilterOrders(c *fiber.Ctx) error {
// 	return nil
// }
