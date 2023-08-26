package routes

import (
	"github.com/aditansh/balkan-task/controllers"
	"github.com/gofiber/fiber/v2"
)

func CartRoutes(app *fiber.App) {

	cart := app.Group("/cart", middleware.verifyUserToken)
	cart.Post("/add", controllers.AddToCart)
	cart.Post("/remove", controllers.RemoveFromCart)
	cart.Post("/update", controllers.UpdateCart)
	cart.Get("/clear", controllers.ClearCart)
	cart.Get("/get", controllers.GetCart)
	cart.Post("/checkout", controllers.Checkout)
}
