package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	UserRoutes(app)
	AdminRoutes(app)
	CartRoutes(app)
	BookRoutes(app)
}
