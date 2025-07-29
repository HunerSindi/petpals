
package routes

import (
	"petapp/internal/handlers"
	"petapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAddressRoutes(app *fiber.App) {
	addressGroup := app.Group("/api/addresses")
	addressGroup.Use(middleware.AuthRequired)

	addressGroup.Post("/", handlers.CreateAddress)
	addressGroup.Get("/", handlers.GetAddresses)
	addressGroup.Put("/:id", handlers.UpdateAddress)
	addressGroup.Delete("/:id", handlers.DeleteAddress)
}
