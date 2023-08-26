package routes

import (
	"github.com/aditansh/balkan-task/controllers"
	"github.com/aditansh/balkan-task/middleware"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {

	auth := app.Group("/auth")

	auth.Post("/verifyotp", controllers.VerifyOTP)
	auth.Post("/resendotp", controllers.ResendOTP)
	auth.Post("/forgot", controllers.ForgotPasswordMail)
	auth.Post("/reset", controllers.ResetPassword)
	auth.Post("/updatepassword", controllers.UpdatePassword)
	auth.Post("/refresh", controllers.RefreshToken)

	auth.Post("/logout", middleware.VerifyToken, controllers.Logout)
	auth.Get("/deactivate", middleware.VerifyToken, controllers.DeactivateAccount)
	auth.Get("/delete", middleware.VerifyToken, controllers.DeleteAccount)
}
