package routes

import (
	"github.com/aditansh/balkan-task/controllers"
	"github.com/aditansh/balkan-task/middleware"
	"github.com/gofiber/fiber/v2"
)

func AdminRoutes(app *fiber.App) {

	app.Post("/admin/register", controllers.RegisterAdmin)

	admin := app.Group("/admin", middleware.VerifyAdminToken)
	admin.Post("/update", controllers.UpdateUser)
	admin.Get("/me", controllers.GetAdminProfile)

	//modify users
	admin.Get("/users", controllers.GetAllUsers)
	admin.Get("/user/:username", controllers.GetUser)
	admin.Post("/user/register", controllers.RegisterUser)
	admin.Post("/user/update/:username", controllers.UpdateUser)
	admin.Post("/user/promote", controllers.MakeAdmins)
	admin.Post("/user/deactivate", controllers.DeactivateUsers)
	admin.Post("/user/delete", controllers.DeleteUsers)

	//modify user carts
	admin.Get("/carts", controllers.GetAllCarts)
	admin.Post("/user/cart", controllers.GetUserCart)
	admin.Post("/user/cart/modify", controllers.ModifyUserCart)
	admin.Post("/user/cart/clear", controllers.ClearUserCart)

	//orders
	admin.Get("/orders", controllers.GetAllOrders)
	admin.Post("orders/search", controllers.SearchOrders)
	admin.Post("orders/filter", controllers.FilterOrders)
	admin.Post("/user/orders", controllers.GetUserOrders)
	admin.Post("/user/order", controllers.GetUserOrder)
	admin.Post("/user/orders/search", controllers.SearchOrders)
	admin.Post("/user/orders/filter", controllers.FilterOrders)

	//books
	admin.Get("/books", controllers.GetBooks)
	admin.Get("/book/:id", controllers.GetBook)
	admin.Post("/book/create", controllers.CreateBook)
	admin.Post("/book/update", controllers.UpdateBook)
	admin.Post("/book/delete", controllers.DeleteBooks)

	//reviews
	admin.Get("/reviews", controllers.GetAllReviews)
	admin.Post("/book/reviews", controllers.GetReviews)
	admin.Get("/book/review/:id", controllers.GetReview)
	admin.Post("/book/review/Modify", controllers.UpdateReview)
	admin.Post("/book/review/delete", controllers.DeleteReview)

	//vendors
	admin.Get("/vendors", controllers.GetAllVendors)
	admin.Get("/vendor/:email", controllers.GetVendor)
	admin.Post("/vendor/register", controllers.RegisterVendor)
	admin.Post("/vendor/update/:email", controllers.UpdateVendor)
	admin.Post("/vendor/deactivate", controllers.DeactivateVendors)
	admin.Post("/vendor/delete", controllers.DeleteVendors)
	admin.Post("/vendor/approve", controllers.ApproveVendor)

	//moderation
	admin.Post("/flaguser", controllers.FlagUser)
	admin.Post("/flagvendor", controllers.FlagVendor)
	admin.Get("/flag/users", controllers.GetFlaggedUsers)
	admin.Get("/flag/vendors", controllers.GetFlaggedVendors)
}
