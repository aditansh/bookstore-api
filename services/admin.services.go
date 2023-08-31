package services

import (
	"fmt"
	"time"

	"github.com/aditansh/balkan-task/database"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
)

//users

func MakeAdmins(usernames []string) *fiber.Error {
	for _, username := range usernames {
		fmt.Println(username)
		user, err := GetUserByUsername(username)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		errr := UpdateRole(&user, "admin")
		if errr != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func DeactivateUsers(usernames []string) *fiber.Error {
	for _, username := range usernames {
		user, err := GetUserByUsername(username)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		updates := make(map[string]interface{})
		updates["is_active"] = false
		updates["updated_at"] = time.Now()

		result := database.DB.Model(&user).Updates(updates)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
		}
	}

	return nil
}

func DeleteUsers(usernames []string) *fiber.Error {
	for _, username := range usernames {
		user, err := GetUserByUsername(username)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		updates := make(map[string]interface{})
		updates["is_active"] = false
		updates["is_deleted"] = true
		updates["updated_at"] = time.Now()

		result := database.DB.Model(&user).Updates(updates)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
		}
	}

	return nil
}

func FlagUser(email string) *fiber.Error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if user.IsFlagged {
		updates := make(map[string]interface{})
		updates["is_deleted"] = true
		updates["is_active"] = false
		updates["updated_at"] = time.Now()

		result := database.DB.Model(&user).Updates(updates)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
		}

		body := "Your account has been deleted due to multiple violations of our policy. If you feel this is an error, please reach out within 5 days."
		err = utils.SendEmail(user.Email, "Policy Violation", body)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
		}
	} else {
		updates := make(map[string]interface{})
		updates["is_flagged"] = true
		updates["updated_at"] = time.Now()

		result := database.DB.Model(&user).Updates(updates)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
		}

		body := "Your account has been flagged due to multiple violations of our policy. The second offence will lead to your account being deleted."
		err = utils.SendEmail(user.Email, "Policy Violation", body)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
		}
	}

	return nil
}

// vendors
func DeactivateVendors(emails []string) *fiber.Error {
	for _, email := range emails {
		vendor, err := GetVendorByEmail(email)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		updates := make(map[string]interface{})
		updates["is_active"] = false
		updates["updated_at"] = time.Now()

		result := database.DB.Model(&vendor).Updates(updates)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error updating vendor")
		}
	}

	return nil
}

func DeleteVendors(emails []string) *fiber.Error {
	for _, email := range emails {
		vendor, err := GetVendorByEmail(email)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		updates := make(map[string]interface{})
		updates["is_active"] = false
		updates["is_deleted"] = true
		updates["updated_at"] = time.Now()

		result := database.DB.Model(&vendor).Updates(updates)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error updating vendor")
		}
	}

	return nil
}

func FlagVendor(email string) *fiber.Error {
	vendor, err := GetVendorByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if vendor.IsFlagged {
		updates := make(map[string]interface{})
		updates["is_deleted"] = true
		updates["is_active"] = false
		updates["updated_at"] = time.Now()

		result := database.DB.Model(&vendor).Updates(updates)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error updating vendor")
		}

		body := "Your account has been deleted due to multiple violations of our policy. If you feel this is an error, please reach out within 5 days."
		err = utils.SendEmail(vendor.Email, "Policy Violation", body)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
		}
	} else {
		updates := make(map[string]interface{})
		updates["is_flagged"] = true
		updates["updated_at"] = time.Now()

		result := database.DB.Model(&vendor).Updates(updates)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error updating vendor")
		}

		body := "Your account has been flagged due to multiple violations of our policy. The second offence will lead to your account being deleted."
		err = utils.SendEmail(vendor.Email, "Policy Violation", body)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
		}
	}

	return nil
}

func ApproveVendor(email string) *fiber.Error {
	vendor, err := GetVendorByEmail(email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if !vendor.IsVerified {
		return fiber.NewError(fiber.StatusBadRequest, "Vendor not verified")
	}

	updates := make(map[string]interface{})
	updates["is_approved"] = true
	updates["is_flagged"] = false
	updates["is_active"] = true
	updates["is_deleted"] = false
	updates["updated_at"] = time.Now()

	result := database.DB.Model(&vendor).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating vendor")
	}

	body := "Your account has been approved. You can now start selling your books."
	err = utils.SendEmail(vendor.Email, "Account Approved", body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error sending mail")
	}

	return nil
}
