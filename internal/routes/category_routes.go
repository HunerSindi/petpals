
package routes

import (
	"petapp/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupCategoryRoutes(app *fiber.App) {
	categoryGroup := app.Group("/api/categories")
	categoryGroup.Get("/", handlers.GetCategories)
	categoryGroup.Get("/:id/products", handlers.GetProductsByCategory)
}
