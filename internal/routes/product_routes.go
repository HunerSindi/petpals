
package routes

import (
	"petapp/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(app *fiber.App) {
	productGroup := app.Group("/api/products")
	productGroup.Get("/:id", handlers.GetProduct)
	productGroup.Get("/:id/images", handlers.GetProductImages)
}
