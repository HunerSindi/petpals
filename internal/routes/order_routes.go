
package routes

import (
	"petapp/internal/handlers"
	"petapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupOrderRoutes(app *fiber.App) {
	orderGroup := app.Group("/api/orders")
	orderGroup.Use(middleware.AuthRequired)

	orderGroup.Post("/", handlers.CreateOrder)
	orderGroup.Get("/", handlers.GetOrders)
	orderGroup.Get("/:id", handlers.GetOrderDetails)
	orderGroup.Put("/:id/cancel", handlers.CancelOrder)
}
