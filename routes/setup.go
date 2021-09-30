package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satriyoaji/sagara-backend-test/controllers"
)

func Setup(app *fiber.App)  {

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)

	app.Get("/products", controllers.FindAllProduct)
	app.Get("/products/:id", controllers.FindProductById)
	app.Post("/products", controllers.CreateProduct)
	app.Patch("/products/:id", controllers.UpdateProductById)
	app.Delete("/products/:id", controllers.DeleteProductById)
}
