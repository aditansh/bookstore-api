package services

import (
	"fmt"

	"github.com/aditansh/balkan-task/database"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
)

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
		updates["is_deleted"] = true

		result := database.DB.Model(&user).Updates(updates)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error updating user")
		}
	}

	return nil
}

func FlagUser(username string) *fiber.Error {
	user, err := GetUserByUsername(username)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if user.IsFlagged {
		updates := make(map[string]interface{})
		updates["is_deleted"] = true
		updates["is_active"] = false

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
