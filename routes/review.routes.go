package routes

import (
	"github.com/aditansh/balkan-task/controllers"
	"github.com/aditansh/balkan-task/middleware"
	"github.com/gofiber/fiber/v2"
)

func ReviewRoutes(app *fiber.App) {

	review := app.Group("/review", middleware.VerifyUserToken)
	review.Get("/get", controllers.GetReviews)
	review.Get("/get/:id", controllers.GetReview)
	review.Post("/review/add", controllers.AddReview)
	review.Post("/review/update", controllers.UpdateReview)
	review.Post("/review/delete", controllers.DeleteReview)
}
