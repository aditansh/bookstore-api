package controllers

import (
	"github.com/aditansh/balkan-task/schemas"
	"github.com/aditansh/balkan-task/services"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterVendor(c *fiber.Ctx) error {
	var payload schemas.RegisterVendorSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	errors := utils.ValidateStruct(&payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"errors": errors,
		})
	}

	err := services.RegisterVendor(&payload)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Vendor registered successfully",
	})
}

func LoginVendor(c *fiber.Ctx) error {
	var payload schemas.LoginEmailSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"errors": errors,
		})
	}

	token, err := services.LoginVendor(&payload)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Vendor logged in successfully",
		"token":   token,
	})
}

func VerifyVendorOTP(c *fiber.Ctx) error {
	var payload schemas.VerifyOTPSchema

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

	email := payload.Email
	err := services.VerifyVendorOTP(&payload, email)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "OTP verified successfully. Please wait for your profile to be approved",
	})
}

func ResendVendorOTP(c *fiber.Ctx) error {
	var payload *schemas.ResendOTPSchema

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

	err := services.ResendVendorOTP(payload.Email)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "OTP resent successfully",
	})
}

func ForgotVendorPasswordMail(c *fiber.Ctx) error {
	var payload *schemas.ForgotPasswordSchema

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

	err := services.SendVendorOTP(payload.Email)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "OTP sent successfully",
	})
}

func ResetVendorPassword(c *fiber.Ctx) error {
	var payload schemas.ResetPasswordSchema

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

	err := services.ResetVendorPassword(&payload)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Password reset successfully",
	})
}

func UpdateVendorPassword(c *fiber.Ctx) error {
	var payload schemas.UpdatePasswordSchema

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

	err := services.UpdateVendorPassword(&payload, ID)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Password updated successfully",
	})
}

func RefreshVendorToken(c *fiber.Ctx) error {
	var payload schemas.RefreshTokenSchema

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

	token, err := services.RefreshVendorAccessToken(&payload)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"token":  token,
	})
}

func UpdateVendor(c *fiber.Ctx) error {
	var payload schemas.UpdateVendorSchema
	var email = c.Params("email")
	var ID *uuid.UUID

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

	if email == "" {
		temp := c.Locals("ID").(uuid.UUID)
		ID = &temp
	} else {
		vendor, err := services.GetVendorByEmail(email)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}
		ID = &vendor.ID
	}

	err := services.UpdateVendor(&payload, *ID)
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

func GetVendorProfile(c *fiber.Ctx) error {
	ID := c.Locals("ID").(uuid.UUID)

	vendor, err := services.GetVendorByID(ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"vendor": vendor,
	})
}

func LogoutVendor(c *fiber.Ctx) error {
	var payload schemas.LogoutSchema

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

	err := services.Logout(&payload)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Logged out successfully",
	})
}

func DeactivateVendorAccount(c *fiber.Ctx) error {
	ID := c.Locals("ID").(uuid.UUID)

	err := services.DeactivateVendorAccount(ID)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Account deactivated successfully",
	})
}

func DeleteVendorAccount(c *fiber.Ctx) error {
	ID := c.Locals("ID").(uuid.UUID)

	err := services.DeleteVendorAccount(ID)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Account deleted successfully",
	})
}

// admin only routes

func GetAllVendors(c *fiber.Ctx) error {
	vendors, err := services.GetAllVendors()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"vendors": vendors,
	})
}

func GetVendor(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Email not provided",
		})
	}

	vendor, err := services.GetVendorByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"vendor": vendor,
	})
}

func DeactivateVendors(c *fiber.Ctx) error {
	var payload schemas.DeactivateDeleteVendorsSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request",
		})
	}

	errors := utils.ValidateStruct(&payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"errors": errors,
		})
	}

	err := services.DeactivateVendors(payload.Emails)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Vendors deactivated successfully",
	})
}

func DeleteVendors(c *fiber.Ctx) error {
	var payload schemas.DeactivateDeleteVendorsSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request",
		})
	}

	errors := utils.ValidateStruct(&payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"errors": errors,
		})
	}

	err := services.DeleteVendors(payload.Emails)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Vendors deleted successfully",
	})
}

func ApproveVendor(c *fiber.Ctx) error {
	var payload schemas.ApproveVendorSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request",
		})
	}

	errors := utils.ValidateStruct(&payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"errors": errors,
		})
	}

	err := services.ApproveVendor(payload.Email)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"status":  false,
			"message": err.Message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Vendor approved successfully",
	})
}
