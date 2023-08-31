package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/aditansh/balkan-task/database"
	"github.com/aditansh/balkan-task/models"
	"github.com/aditansh/balkan-task/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllCarts() ([]models.Cart, error) {
	var carts []models.Cart
	result := database.DB.Preload("Items").Find(&carts)
	if result.Error != nil {
		return nil, result.Error
	}

	return carts, nil
}

// GetCart returns the cart of the user
func GetCart(ID uuid.UUID) (models.Cart, error) {
	var cart models.Cart
	result := database.DB.Preload("Items").Where("user_id = ?", ID).First(&cart)
	if result.Error != nil {
		return models.Cart{}, result.Error
	}

	return cart, nil
}

func GetCartItem(cartID uuid.UUID, bookID uuid.UUID) (models.CartItem, error) {
	var cartItem models.CartItem
	result := database.DB.Where("cart_id = ? AND book_id = ?", cartID, bookID).First(&cartItem)
	if result.Error != nil {
		return models.CartItem{}, result.Error
	}

	return cartItem, nil
}

func AddToCart(ID uuid.UUID, payload *schemas.ModifyCartSchema) *fiber.Error {
	cart, err := GetUserCart(ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bookID, err := uuid.Parse(payload.BookID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid book id")
	}

	book, err := GetBookByID(bookID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fmt.Println(book.Name)

	newCartItem := models.CartItem{
		CartID:   cart.ID,
		BookID:   bookID,
		BookName: book.Name,
		Quantity: payload.Quantity,
		VendorID: book.VendorID,
	}

	result := database.DB.Create(&newCartItem)
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint \"cart_items_pkey\"") {
		return fiber.NewError(fiber.StatusConflict, "Book exists in cart")
	} else if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func UpdateCart(ID uuid.UUID, payload *schemas.ModifyCartSchema) *fiber.Error {
	cart, err := GetUserCart(ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bookID, err := uuid.Parse(payload.BookID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid book id")
	}

	cartItem, err := GetCartItem(cart.ID, bookID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	updates := make(map[string]interface{})
	updates["quantity"] = payload.Quantity
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&cartItem).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func ClearCart(ID uuid.UUID) *fiber.Error {
	cart, err := GetUserCart(ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var cartItem models.CartItem
	result := database.DB.Where("cart_id = ?", cart.ID).Delete(&cartItem)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func RemoveFromCart(ID uuid.UUID, payload *schemas.RemoveFromCart) *fiber.Error {
	cart, err := GetUserCart(ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bookID, err := uuid.Parse(payload.BookID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid book id")
	}

	fmt.Println(cart.ID, bookID)
	var cartItem models.CartItem
	result := database.DB.Where("cart_id = ? AND book_id = ?", cart.ID, bookID).Delete(&cartItem)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func Checkout(ID uuid.UUID, payload *schemas.CheckoutSchema) *fiber.Error {
	cart, err := GetUserCart(ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	order, err := CreateOrder(ID, payload.Address)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var orderItems []models.OrderItem
	var totalCost = 0.0

	for _, cartItem := range cart.Items {
		book, err := GetBookByID(cartItem.BookID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		vendor, err := GetVendorByID(book.VendorID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		orderItem := models.OrderItem{
			OrderID:    order.ID,
			BookName:   book.Name,
			Author:     book.Author,
			Quantity:   cartItem.Quantity,
			VendorName: vendor.Name,
		}

		totalCost += float64(cartItem.Quantity) * book.Cost

		orderItems = append(orderItems, orderItem)
	}

	updates := make(map[string]interface{})
	updates["value"] = totalCost
	updates["items"] = orderItems
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&order).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	var cartItem models.CartItem
	result = database.DB.Where("cart_id = ?", cart.ID).Delete(&cartItem)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}
