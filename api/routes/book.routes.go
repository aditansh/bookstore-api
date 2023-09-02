package routes

import (
	"github.com/aditansh/balkan-task/controllers"
	"github.com/aditansh/balkan-task/middleware"
	"github.com/gofiber/fiber/v2"
)

func BookRoutes(app *fiber.App) {

	book := app.Group("/book", middleware.VerifyUserToken)
	book.Get("/get/:id", controllers.GetBook)
	book.Get("/getall", controllers.GetBooks)
	book.Post("/search", controllers.SearchBooks)
	book.Post("/filter", controllers.FilterBooks)
}

func BookAdminRoutes(app *fiber.App) {

	bookAdmin := app.Group("/bookadmin", middleware.VerifyVendorToken)
	bookAdmin.Post("/create", controllers.CreateBook)
	bookAdmin.Post("/update/:id", controllers.UpdateBook)
	bookAdmin.Post("/delete", controllers.DeleteBooks)
}
