package routes

import (
	"github.com/aditansh/balkan-task/controllers"
	"github.com/aditansh/balkan-task/middleware"
	"github.com/gofiber/fiber/v2"
)

func AdminRoutes(app *fiber.App) {

	app.Post("/register", controllers.RegisterAdmin)
	app.Post("/login", controllers.LoginAdmin)

	admin := app.Group("/admin", middleware.VerifyAdminToken)
	admin.Post("/update", controllers.UpdateAdmin)
	admin.Get("/me", controllers.GetAdminProfile)

	//modify users
	admin.Get("/users", controllers.GetAllUsers)
	admin.Get("/user/:id", controllers.GetUser)
	admin.Post("/user/register", controllers.RegisterUser)
	admin.Post("/user/update", controllers.UpdateUser)
	admin.Post("/user/deactivate", controllers.DeactivateAccount)
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
	admin.Get("/vendor/:id", controllers.GetVendor)
	admin.Post("/vendor/register", controllers.RegisterVendor)
	admin.Post("/vendor/update", controllers.UpdateVendor)
	admin.Post("/vendor/deactivate", controllers.DeactivateAccount)
	admin.Post("/vendor/delete", controllers.DeleteVendors)
	admin.Post("/vendor/approve", controllers.ApproveVendor)

	//moderation
	admin.Post("/flaguser", controllers.FlagUser)
	admin.Post("/flagvendor", controllers.FlagVendor)
	admin.Get("/flag/user", controllers.GetFlaggedUsers)
	admin.Get("/flag/vendor", controllers.GetFlaggedVendors)
}
