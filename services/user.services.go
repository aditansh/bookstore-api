package services

import (
	"fmt"
	"time"

	"github.com/aditansh/balkan-task/cache"
	"github.com/aditansh/balkan-task/database"
	"github.com/aditansh/balkan-task/models"
	"github.com/aditansh/balkan-task/schemas"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		return []models.User{}, result.Error
	}

	return users, nil
}

func GetUserByID(id uuid.UUID) (models.User, error) {
	var user models.User
	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	if user.IsDeleted && user.IsFlagged {
		return models.User{}, fmt.Errorf("account is deleted")
	}

	return user, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	if user.IsDeleted && user.IsFlagged {
		return models.User{}, fmt.Errorf("account is deleted")
	}

	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	if user.IsDeleted && user.IsFlagged {
		return models.User{}, fmt.Errorf("account is deleted")
	}

	return user, nil
}

func GetFlaggedUsers() ([]models.User, error) {
	var users []models.User
	result := database.DB.Where("is_flagged = ?", true).Find(&users)
	if result.Error != nil {
		return []models.User{}, result.Error
	}

	return users, nil
}

func RegisterUser(payload *schemas.RegisterUserSchema, isAdmin bool) *fiber.Error {

	_, check1 := GetUserByEmail(payload.Email)
	if check1 == nil {
		return fiber.NewError(fiber.StatusConflict, "Email already exists")
	}

	_, check2 := GetUserByUsername(payload.Username)
	if check2 == nil {
		return fiber.NewError(fiber.StatusConflict, "Username already exists")
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error hashing password")
	}

	otp, _ := utils.GenerateOTP(6)
	//clear any otp set for this email due to previous failed attempts
	err = cache.DeleteValue(payload.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}
	err = cache.SetValue(payload.Email, otp, 48*time.Hour)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error storing OTP")
	}

	body := "Your OTP is " + otp + ".\nIt will expire in 48 hours."
	err = utils.SendEmail(payload.Email, "OTP Verification", body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
	}

	var role string
	if isAdmin {
		role = "admin"
	} else {
		role = "user"
	}

	newUser := models.User{
		Name:       payload.Name,
		Username:   payload.Username,
		Email:      payload.Email,
		Password:   hashedPassword,
		IsVerified: false,
		IsDeleted:  false,
		Role:       role,
	}

	result := database.DB.Create(&newUser)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func UpdateUser(payload *schemas.UpdateUserSchema, ID uuid.UUID) *fiber.Error {
	user, err := GetUserByID(ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	updates := make(map[string]interface{})

	if payload.Name != "" {
		updates["name"] = payload.Name
	}

	if payload.Username != "" {
		_, check := GetUserByUsername(payload.Username)
		if check == nil {
			return fiber.NewError(fiber.StatusConflict, "Username already exists")
		}
		updates["username"] = payload.Username
	}

	if payload.Email != "" {
		_, check := GetUserByEmail(payload.Email)
		if check == nil {
			return fiber.NewError(fiber.StatusConflict, "Email already exists")
		}
		updates["email"] = payload.Email
		updates["is_verified"] = false
		otp, _ := utils.GenerateOTP(6)
		//clear any otp set for this email due to previous failed attempts
		err = cache.DeleteValue(payload.Email)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
		}

		err = cache.SetValue(payload.Email, otp, 48*time.Hour)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error storing OTP")
		}

		body := "Your OTP is " + otp + ".\nIt will expire in 48 hours."
		err = utils.SendEmail(payload.Email, "OTP Verification", body)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
		}
	}
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
	}

	return nil
}

func UpdateRole(user *models.User, role string) *fiber.Error {
	updates := make(map[string]interface{})
	updates["role"] = role
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
	}

	return nil
}
