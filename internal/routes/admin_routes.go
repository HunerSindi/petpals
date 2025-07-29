
package routes

import (
	"petapp/internal/handlers"
	"petapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app *fiber.App) {
	adminGroup := app.Group("/api/admin")
	adminGroup.Use(middleware.AdminAuthRequired)

	// Clinic Management
	adminGroup.Post("/clinics", handlers.AdminCreateClinic)
	adminGroup.Get("/clinics", handlers.AdminGetClinics)
	adminGroup.Put("/clinics/:id", handlers.AdminUpdateClinic)
	adminGroup.Delete("/clinics/:id", handlers.AdminDeleteClinic)

	// Product Management
	adminGroup.Post("/products", handlers.AdminCreateProduct)
	adminGroup.Put("/products/:id", handlers.AdminUpdateProduct)
	adminGroup.Delete("/products/:id", handlers.AdminDeleteProduct)
	adminGroup.Post("/products/:id/images", handlers.AdminAddProductImages)

	// Order Management
	adminGroup.Get("/orders", handlers.AdminGetOrders)
	adminGroup.Put("/orders/:id/status", handlers.AdminUpdateOrderStatus)
}
