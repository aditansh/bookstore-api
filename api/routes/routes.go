package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	AdminRoutes(app)
	AuthRoutes(app)
	BookAdminRoutes(app)
	BookRoutes(app)
	CartRoutes(app)
	OrderRoutes(app)
	ReviewRoutes(app)
	UserRoutes(app)
	VendorRoutes(app)
}
