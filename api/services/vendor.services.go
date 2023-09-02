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

func GetAllVendors() ([]models.Vendor, error) {
	var vendors []models.Vendor
	result := database.DB.Find(&vendors)
	if result.Error != nil {
		return []models.Vendor{}, result.Error
	}

	return vendors, nil
}

func GetVendorByID(id uuid.UUID) (models.Vendor, error) {
	var vendor models.Vendor
	result := database.DB.Where("id = ?", id).First(&vendor)
	if result.Error != nil {
		return models.Vendor{}, result.Error
	}

	if vendor.IsDeleted && vendor.IsFlagged {
		return models.Vendor{}, fmt.Errorf("account is deleted")
	}

	return vendor, nil
}

func GetVendorByEmail(email string) (models.Vendor, error) {
	var vendor models.Vendor
	result := database.DB.Where("email = ?", email).First(&vendor)
	if result.Error != nil {
		return models.Vendor{}, result.Error
	}

	if vendor.IsDeleted && vendor.IsFlagged {
		return models.Vendor{}, fmt.Errorf("account is deleted")
	}

	return vendor, nil
}

func GetFlaggedVendors() ([]models.Vendor, error) {
	var vendors []models.Vendor
	result := database.DB.Where("is_flagged = ?", true).Find(&vendors)
	if result.Error != nil {
		return []models.Vendor{}, result.Error
	}

	return vendors, nil
}

func RegisterVendor(payload *schemas.RegisterVendorSchema) *fiber.Error {

	_, check1 := GetVendorByEmail(payload.Email)
	if check1 == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Email already exists")
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
		return fiber.NewError(fiber.StatusInternalServerError, "Error storing otp")
	}

	body := "Your OTP is " + otp + ".\nIt will expire in 48 hours."
	err = utils.SendEmail(payload.Email, "OTP Verification", body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
	}

	newVendor := models.Vendor{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
	}

	result := database.DB.Create(&newVendor)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func LoginVendor(payload *schemas.LoginEmailSchema) (string, *fiber.Error) {

	vendor, err := GetVendorByEmail(payload.Email)
	if err != nil {
		return "", fiber.NewError(fiber.StatusNotFound, "Email not found")
	}

	if !vendor.IsVerified {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Email not verified")
	}

	if !vendor.IsApproved {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Account not approved")
	}

	if !vendor.IsActive {
		updates := make(map[string]interface{})
		updates["is_active"] = true
		updates["updated_at"] = time.Now()

		result := database.DB.Model(&vendor).Updates(updates)
		if result.Error != nil {
			return "", fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
		}
	}

	if vendor.IsDeleted {
		updates := make(map[string]interface{})
		updates["is_deleted"] = false
		updates["updated_at"] = time.Now()

		result := database.DB.Model(&vendor).Updates(updates)
		if result.Error != nil {
			return "", fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
		}
	}

	token, err := utils.GenerateRefreshToken(vendor.ID, vendor.Email)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error generating token")
	}

	err = cache.SetValue(token, vendor.ID.String(), 0)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error storing token")
	}

	return token, nil
}

func VerifyVendorOTP(payload *schemas.VerifyOTPSchema, email string) *fiber.Error {
	vendor, err := GetVendorByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if vendor.IsVerified {
		return fiber.NewError(fiber.StatusConflict, "Vendor is verifed")
	}

	otp, err := cache.GetValue(email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error retrieving OTP")
	}

	if otp != payload.OTP {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid OTP")
	}

	err = cache.DeleteValue(email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete otp")
	}

	updates := make(map[string]interface{})
	updates["is_verified"] = true
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&vendor).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating vendor")
	}

	return nil
}

func ResendVendorOTP(email string) *fiber.Error {
	vendor, err := GetVendorByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if vendor.IsVerified {
		return fiber.NewError(fiber.StatusConflict, "Vendor is verifed")
	}

	otp, _ := utils.GenerateOTP(6)
	err = cache.DeleteValue(email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}
	err = cache.SetValue(email, otp, 48*time.Hour)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error storing OTP")
	}

	body := "Your OTP is " + otp + ".\nIt will expire in 48 hours."
	err = utils.SendEmail(email, "OTP Verification", body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
	}

	return nil
}

func SendVendorOTP(email string) *fiber.Error {
	vendor, err := GetVendorByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if !vendor.IsVerified {
		return fiber.NewError(fiber.StatusConflict, "Vendor is not verified")
	}

	otp, _ := utils.GenerateOTP(6)
	err = cache.DeleteValue(email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}

	err = cache.SetValue(otp, email, time.Hour)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error storing OTP")
	}

	body := "Your OTP to reset password is " + otp + ".\nIt will expire in 1 hour."
	err = utils.SendEmail(email, "Reset Password", body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
	}

	return nil
}

func ResetVendorPassword(payload *schemas.ResetPasswordSchema) *fiber.Error {

	email, err := cache.GetValue(payload.OTP)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "OTP not found")
	}

	vendor, err := GetVendorByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if utils.CheckPasswordHash(payload.Password, vendor.Password) {
		return fiber.NewError(fiber.StatusConflict, "New password cannot be same as old password")
	}

	err = cache.DeleteValue(payload.OTP)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error deleting otp")
	}

	newPass, err := utils.HashPassword(payload.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error hashing password")
	}

	updates := make(map[string]interface{})
	updates["password"] = newPass
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&vendor).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating vendor")
	}

	return nil
}

func UpdateVendorPassword(payload *schemas.UpdatePasswordSchema, ID uuid.UUID) *fiber.Error {
	vendor, err := GetVendorByID(ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if !utils.CheckPasswordHash(payload.OldPassword, vendor.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "Incorrect password")
	}

	if payload.OldPassword == payload.NewPassword {
		return fiber.NewError(fiber.StatusBadRequest, "New Password cannot be the same as old password")
	}

	newPass, err := utils.HashPassword(payload.NewPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error hashing password")
	}

	updates := make(map[string]interface{})
	updates["password"] = newPass
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&vendor).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating vendor")
	}

	return nil
}

func RefreshVendorAccessToken(payload *schemas.RefreshTokenSchema) (string, *fiber.Error) {
	ID, err := cache.GetValue(payload.RefreshToken)
	if err != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	vendorID, err := uuid.Parse(ID)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error parsing vendor id")
	}

	vendor, err := GetVendorByID(vendorID)
	if err != nil {
		return "", fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	token, err := utils.GenerateAccessToken(vendor.ID, vendor.Email)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error generating token")
	}

	return token, nil
}

func UpdateVendor(payload *schemas.UpdateVendorSchema, ID uuid.UUID) *fiber.Error {
	vendor, err := GetVendorByID(ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	updates := make(map[string]interface{})

	if payload.Name != "" {
		updates["name"] = payload.Name
	}

	if payload.Email != "" {
		_, check := GetVendorByEmail(payload.Email)
		if check == nil {
			return fiber.NewError(fiber.StatusBadRequest, "Email already exists")
		}
		updates["email"] = payload.Email
		updates["is_verified"] = false
		otp, _ := utils.GenerateOTP(6)
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

	result := database.DB.Model(&vendor).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating vendor")
	}

	return nil
}

func DeactivateVendorAccount(ID uuid.UUID) *fiber.Error {
	vendor, err := GetVendorByID(ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if !vendor.IsActive {
		return fiber.NewError(fiber.StatusConflict, "Account is already deactivated")
	}

	updates := make(map[string]interface{})
	updates["is_active"] = false
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&vendor).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating vendor")
	}

	return nil
}

func DeleteVendorAccount(ID uuid.UUID) *fiber.Error {
	vendor, err := GetVendorByID(ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if vendor.IsDeleted {
		return fiber.NewError(fiber.StatusConflict, "Account is already deleted")
	}

	updates := make(map[string]interface{})
	updates["is_deleted"] = true
	updates["is_active"] = false
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&vendor).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating vendor")
	}

	body := "Thank you for using our platform. Your account will be deleted in 5 days. If you wish to restore your account, login using the same credentials."
	err = utils.SendEmail(vendor.Email, "Account Deleted", body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
	}

	return nil
}
