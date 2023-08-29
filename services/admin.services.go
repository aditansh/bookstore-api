package services

import (
	"fmt"

	"github.com/aditansh/balkan-task/database"
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
