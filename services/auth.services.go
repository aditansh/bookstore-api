package services

import (
	"time"

	"github.com/aditansh/balkan-task/cache"
	"github.com/aditansh/balkan-task/database"
	"github.com/aditansh/balkan-task/schemas"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func LoginEmail(payload *schemas.LoginEmailSchema) (string, *fiber.Error) {

	user, err := GetUserByEmail(payload.Email)
	if err != nil {
		return "", fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if !user.IsVerified {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Email not verified")
	}

	if !utils.CheckPasswordHash(payload.Password, user.Password) {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Incorrect password")
	}

	if !user.IsActive {
		updates := make(map[string]interface{})
		updates["is_active"] = true

		result := database.DB.Model(&user).Updates(updates)
		if result.Error != nil {
			return "", fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
		}
	}

	if user.IsDeleted {
		updates := make(map[string]interface{})
		updates["is_deleted"] = false

		result := database.DB.Model(&user).Updates(updates)
		if result.Error != nil {
			return "", fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
		}
	}

	token, err := utils.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error generating token")
	}

	err = cache.SetValue(token, user.ID.String(), 0)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error storing token")
	}

	return token, nil
}

func LoginUsername(payload *schemas.LoginUsernameSchema) (string, *fiber.Error) {

	user, err := GetUserByUsername(payload.Username)
	if err != nil {
		return "", fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if !user.IsVerified {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Email not verified")
	}

	if !utils.CheckPasswordHash(payload.Password, user.Password) {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Incorrect password")
	}

	if !user.IsActive {
		updates := make(map[string]interface{})
		updates["is_active"] = true

		result := database.DB.Model(&user).Updates(updates)
		if result.Error != nil {
			return "", fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
		}
	}

	if user.IsDeleted {
		updates := make(map[string]interface{})
		updates["is_deleted"] = false

		result := database.DB.Model(&user).Updates(updates)
		if result.Error != nil {
			return "", fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
		}
	}

	token, err := utils.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error generating token")
	}

	err = cache.SetValue(token, user.ID.String(), 0)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error storing token")
	}

	return token, nil
}

func VerifyOTP(payload *schemas.VerifyOTPSchema, email string) *fiber.Error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if user.IsVerified {
		return fiber.NewError(fiber.StatusConflict, "User is verifed")
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

	result := database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
	}

	return nil
}

func ResendOTP(email string) *fiber.Error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if user.IsVerified {
		return fiber.NewError(fiber.StatusConflict, "User is verifed")
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

func SendOTP(email string) *fiber.Error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if !user.IsVerified {
		return fiber.NewError(fiber.StatusConflict, "User is not verified")
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

func ResetPassword(payload *schemas.ResetPasswordSchema) *fiber.Error {

	email, err := cache.GetValue(payload.OTP)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "OTP not found")
	}

	user, err := GetUserByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if utils.CheckPasswordHash(payload.Password, user.Password) {
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

	result := database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
	}

	return nil
}

func UpdatePassword(payload *schemas.UpdatePasswordSchema, email string) *fiber.Error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if !utils.CheckPasswordHash(payload.OldPassword, user.Password) {
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

	result := database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
	}

	return nil
}

func RefreshAccessToken(payload *schemas.RefreshTokenSchema) (string, *fiber.Error) {
	ID, err := cache.GetValue(payload.RefreshToken)
	if err != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	userID, err := uuid.Parse(ID)
	if err != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Error parsing user id")
	}

	user, err := GetUserByID(userID)
	if err != nil {
		return "", fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	token, err := utils.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Error generating token")
	}

	return token, nil
}

func Logout(payload *schemas.LogoutSchema) *fiber.Error {
	id, err := cache.GetValue(payload.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	_, err = uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	err = cache.DeleteValue(payload.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	return nil
}

func DeactivateAccount(email string) *fiber.Error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if !user.IsActive {
		return fiber.NewError(fiber.StatusConflict, "Account already deactivated")
	}

	updates := make(map[string]interface{})
	updates["is_active"] = false

	result := database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
	}

	return nil
}

func DeleteAccount(email string) *fiber.Error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if user.IsDeleted {
		return fiber.NewError(fiber.StatusConflict, "Account already deleted")
	}

	updates := make(map[string]interface{})
	updates["is_deleted"] = true
	updates["is_active"] = false

	result := database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
	}

	body := "Thank you for using our platform. Your account will be deleted in 5 days. If you wish to restore your account, login using the same credentials."
	err = utils.SendEmail(email, "Account Deleted", body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
	}

	return nil
}
