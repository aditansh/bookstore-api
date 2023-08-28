package routes

import (
	"github.com/aditansh/balkan-task/controllers"
	"github.com/aditansh/balkan-task/middleware"
	"github.com/gofiber/fiber/v2"
)

func VendorRoutes(app *fiber.App) {

	vendor := app.Group("/vendor")
	vendor.Post("/register", controllers.RegisterVendor)
	vendor.Post("/login", controllers.LoginVendor)
	vendor.Post("/verifyotp", controllers.VerifyVendorOTP)
	vendor.Post("/resendotp", controllers.ResendVendorOTP)
	vendor.Post("/forgot", controllers.ForgotVendorPasswordMail)
	vendor.Post("/refresh", controllers.RefreshVendorToken)

	vendor.Post("/reset", middleware.VerifyVendorToken, controllers.ResetVendorPassword)
	vendor.Post("/updatepassword", middleware.VerifyVendorToken, controllers.UpdateVendorPassword)
	vendor.Post("/update", middleware.VerifyVendorToken, controllers.UpdateVendor)
	vendor.Get("/me", middleware.VerifyVendorToken, controllers.GetVendorProfile)
}
