package routes

import (
	"github.com/aditansh/balkan-task/controllers"
	"github.com/aditansh/balkan-task/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {

	user := app.Group("/user")
	user.Post("/register", controllers.RegisterUser)

	user.Post("/update", middleware.VerifyUserToken, controllers.UpdateUser)
	user.Get("/me", middleware.VerifyUserToken, controllers.GetUserProfile)
}
