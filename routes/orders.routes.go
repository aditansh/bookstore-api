package routes

import (
	"github.com/aditansh/balkan-task/controllers"
	"github.com/aditansh/balkan-task/middleware"
	"github.com/gofiber/fiber/v2"
)

func OrderRoutes(app *fiber.App) {

	orders := app.Group("/orders", middleware.VerifyUserToken)
	orders.Get("/get", controllers.GetOrders)
	orders.Get("/get/:id", controllers.GetOrder)
	orders.Post("/search", controllers.SearchOrders)
	orders.Post("/filter", controllers.FilterOrders)
}
