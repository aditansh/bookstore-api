package services

import (
	"github.com/aditansh/balkan-task/database"
	"github.com/aditansh/balkan-task/models"
	"github.com/aditansh/balkan-task/schemas"
	"github.com/aditansh/balkan-task/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllReviews() ([]models.Review, error) {
	var reviews []models.Review
	result := database.DB.Find(&reviews)
	if result.Error != nil {
		return []models.Review{}, result.Error
	}

	return reviews, nil
}

func GetReview(ID uuid.UUID) (models.Review, error) {
	var review models.Review
	result := database.DB.Where("id = ?", ID).First(&review)
	if result.Error != nil {
		return models.Review{}, result.Error
	}

	return review, nil
}

func AddReview(payload *schemas.AddReviewSchema, userID uuid.UUID) *fiber.Error {
	bookID, err := uuid.Parse(payload.BookID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user, err := GetUserByID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if check := utils.StringInSlice(payload.BookID, user.Books); !check {
		return fiber.NewError(fiber.StatusBadRequest, "You can only review books that you have bought")
	}

	var review = models.Review{
		UserID:   userID,
		BookID:   bookID,
		Username: user.Username,
		Comment:  payload.Comment,
		Rating:   payload.Rating,
	}

	result := database.DB.Create(&review)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func UpdateReview(payload *schemas.UpdateReviewSchema, userID uuid.UUID) *fiber.Error {
	reviewID, err := uuid.Parse(payload.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	review, err := GetReview(reviewID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	updates := make(map[string]interface{})
	updates["comment"] = payload.Comment
	updates["rating"] = payload.Rating

	result := database.DB.Model(&review).Updates(updates)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func DeleteReview(payload *schemas.DeleteReviewSchema, userID uuid.UUID) *fiber.Error {
	reviewID, err := uuid.Parse(payload.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	review, err := GetReview(reviewID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user, err := GetUserByID(userID)

	if review.UserID != userID && user.Role == "user" {
		return fiber.NewError(fiber.StatusUnauthorized, "You can only delete your own reviews")
	}

	result := database.DB.Delete(&review)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}

func DeleteReviews(payload *schemas.DeleteReviewsSchema) *fiber.Error {
	var review models.Review
	for _, id := range payload.IDs {
		reviewID, err := uuid.Parse(id)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid review ID")
		}

		result := database.DB.Where("id = ?", reviewID).Delete(&review)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
		}
	}

	return nil
}
