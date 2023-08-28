package routes

import (
	"github.com/aditansh/balkan-task/controllers"
	"github.com/aditansh/balkan-task/middleware"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {

	auth := app.Group("/auth")

	auth.Post("/login", controllers.Login)
	auth.Post("/verifyotp", controllers.VerifyOTP)
	auth.Post("/resendotp", controllers.ResendOTP)
	auth.Post("/forgot", controllers.ForgotPasswordMail)
	auth.Post("/reset", controllers.ResetPassword)
	auth.Post("/refresh", controllers.RefreshToken)

	auth.Post("/updatepassword", middleware.VerifyToken, controllers.UpdatePassword)
	auth.Post("/logout", middleware.VerifyToken, controllers.Logout)
	auth.Get("/deactivate", middleware.VerifyToken, controllers.DeactivateAccount)
	auth.Get("/delete", middleware.VerifyToken, controllers.DeleteAccount)
}
