package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/aditansh/balkan-task/database"
	"github.com/aditansh/balkan-task/models"
	"github.com/aditansh/balkan-task/schemas"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllCarts() ([]models.Cart, error) {
	var carts []models.Cart
	result := database.DB.Preload("Items").Find(&carts)
	if result.Error != nil {
		return []models.Cart{}, result.Error
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

	if payload.Quantity > book.Stock {
		return fiber.NewError(fiber.StatusBadRequest, "not enough stock available")
	}

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

	book, err := GetBookByID(bookID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if payload.Quantity > book.Stock {
		return fiber.NewError(fiber.StatusBadRequest, "not enough stock available")
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

	var cartItem models.CartItem
	result := database.DB.Where("cart_id = ? AND book_id = ?", cart.ID, bookID).Delete(&cartItem)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func Checkout(ID uuid.UUID, payload *schemas.CheckoutSchema) (string, *fiber.Error) {
	user, err := GetUserByID(ID)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	userBooks := user.Books

	cart, err := GetUserCart(ID)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(cart.Items) == 0 {
		return "", fiber.NewError(fiber.StatusBadRequest, "cart is empty")
	}

	order, err := CreateOrder(ID, payload.Address)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, cartItem := range cart.Items {
		book, err := GetBookByID(cartItem.BookID)
		if err != nil {
			return order.ID.String(), fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if cartItem.Quantity > book.Stock {
			return order.ID.String(), fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("not enough stock for %s", book.Name))
		}

		vendor, err := GetVendorByID(book.VendorID)
		if err != nil {
			return order.ID.String(), fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		orderItem := models.OrderItem{
			OrderID:    order.ID,
			BookID:     book.ID,
			BookName:   book.Name,
			Author:     book.Author,
			Quantity:   cartItem.Quantity,
			VendorName: vendor.Name,
		}

		result := database.DB.Create(&orderItem)
		if result.Error != nil {
			return order.ID.String(), fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
		}

		updates := make(map[string]interface{})
		updates["stock"] = book.Stock - cartItem.Quantity
		updates["updated_at"] = time.Now()

		result = database.DB.Model(&book).Updates(updates)
		if result.Error != nil {
			return order.ID.String(), fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
		}

		userBooks = utils.AppendUnique(userBooks, book.ID.String())
	}

	cartValue, err := GetCartValue(ID)
	if err != nil {
		return order.ID.String(), fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	updates := make(map[string]interface{})
	updates["value"] = cartValue
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&order).Updates(updates)
	if result.Error != nil {
		return order.ID.String(), fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	updates = make(map[string]interface{})
	updates["books"] = userBooks

	result = database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return order.ID.String(), fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	errr := ClearCart(ID)
	if errr != nil {
		return "", errr
	}

	return "", nil
}

func GetCartValue(ID uuid.UUID) (float64, error) {
	cart, err := GetUserCart(ID)
	if err != nil {
		return 0, err
	}

	if len(cart.Items) == 0 {
		return 0, fmt.Errorf("no items in cart")
	}

	var value float64
	for _, cartItem := range cart.Items {
		book, err := GetBookByID(cartItem.BookID)
		if err != nil {
			return 0, err
		}

		value += book.Cost * float64(cartItem.Quantity)
	}

	return value, nil
}
