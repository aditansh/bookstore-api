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

	vendor.Post("/update", middleware.VerifyVendorToken, controllers.UpdateVendor)
	vendor.Get("/me", middleware.VerifyVendorToken, controllers.GetVendorProfile)
}
